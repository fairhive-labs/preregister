package main

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/gin-gonic/gin"
)

//go:embed templates
var tfs embed.FS

func setupRouter(app *App) *gin.Engine {
	r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	t := template.Must(template.ParseFS(tfs, "templates/*"))
	r.SetHTMLTemplate(t)
	r.Use(app.cors, app.limit)
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/:path1/:path2/count", app.count)
	r.GET("/:path1/:path2/list", app.list)
	r.POST("/register", app.register)
	r.POST("/activate/:token/:hash", app.activate)
	return r
}

var jwtregexp = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)

func generateSecuredLink(t string) string {
	return fmt.Sprintf("http://fairhive.io/activate/%s", t)
}

func (app *App) register(c *gin.Context) {
	var u data.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := app.jwt.Create(&u, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash := app.jwt.Hash(token)
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		sl := generateSecuredLink(token)
		app.mailer.SendActivationEmail(u.Email, sl, hash)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"hash": hash,
	})
}

func (app *App) activate(c *gin.Context) {
	t := c.Param("token")
	h := c.Param("hash")
	if !jwtregexp.MatchString(t) || app.jwt.Hash(t) != h {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	u, err := app.jwt.Extract(t) // verify + extract
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// @TODO: control if user is already present
	// ra, err := app.db.IsPresent(u.Address)
	// 409 Conflict

	// @TODO: control if sponsor is present
	// rs, err := app.db.IsPresent(u.Sponsor)
	// 400 Bad Request

	err = app.db.Save(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.mailer.SendConfirmationEmail(u.Email)
	}()

	c.JSON(http.StatusCreated, nil)
}

func (app *App) limit(c *gin.Context) {
	ip := c.ClientIP()
	l := app.rl.GetAccess(ip)
	if !l.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "Too Many Requests",
			"ip":    ip,
		})
		return
	}
	c.Next()
}

func (app *App) cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func (app *App) count(c *gin.Context) {
	p1, p2 := c.Param("path1"), c.Param("path2")
	if p1 != app.secpath1 || p2 != app.secpath2 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	cn, err := app.db.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t := 0
	for _, v := range cn {
		t += v
	}
	mime := c.DefaultQuery("mime", "html")
	switch mime {
	case "json":
		c.JSON(http.StatusOK, gin.H{
			"users": cn,
			"total": t,
		})
		return
	case "xml":
		type xmlUser struct {
			Type  string
			Value int
		}
		type Count struct {
			Total int
			Users []xmlUser
		}
		u := []xmlUser{}
		for t, v := range cn {
			u = append(u, xmlUser{t, v})
		}
		sort.Slice(u, func(i, j int) bool {
			return u[i].Type < u[j].Type
		})
		c.XML(http.StatusOK, Count{
			Total: t,
			Users: u,
		})
		return
	default:
		c.HTML(http.StatusOK, "count_template.html", gin.H{
			"users": cn,
			"total": t,
		})
		return
	}
}

func (app *App) list(c *gin.Context) {
	p1, p2 := c.Param("path1"), c.Param("path2")
	if p1 != app.secpath1 || p2 != app.secpath2 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	options := []int{}
	offset := c.Query("offset")
	if offset != "" {
		v, err := strconv.Atoi(offset)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		options = append(options, v)
	}
	if max := c.Query("max"); max != "" {
		v, err := strconv.Atoi(max)
		if err != nil || offset == "" { // offset & max required
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		options = append(options, v)
	}

	users, err := app.db.List(options...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mime := c.DefaultQuery("mime", "json")
	switch mime {
	case "csv":
		b := new(bytes.Buffer)
		w := csv.NewWriter(b)
		err := w.Write([]string{"type", "address", "email", "uuid", "timestamp"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, u := range users {
			l, _ := time.LoadLocation("Europe/Paris")
			err := w.Write([]string{u.Type, u.Address, u.Email, u.UUID, time.UnixMilli(u.Timestamp).In(l).String()})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		w.Flush()
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=users_list_%s.csv", time.Now().Format("20060102-150405")))
		c.Data(http.StatusOK, "text/csv", b.Bytes())
		// c.Writer.Write(b.Bytes())
		return
	default:
		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": len(users),
		})
		return
	}
}
