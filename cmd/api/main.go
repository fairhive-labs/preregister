package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fairhive-labs/preregister/internal/crypto"
	"github.com/fairhive-labs/preregister/internal/crypto/cipher"
	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/fairhive-labs/preregister/internal/limiter"
	"github.com/fairhive-labs/preregister/internal/mailer"
)

type App struct {
	db                 data.DB
	jwt                crypto.Token
	mailer             mailer.Mailer
	wg                 sync.WaitGroup
	rl                 *limiter.RateLimiter
	secpath1, secpath2 string
}

//go:embed templates
var tfs embed.FS

const (
	tableName = "Waitlist"
)

var (
	jwts               = map[string]crypto.Token{}
	ek                 string
	secpath1, secpath2 string
)

func init() {
	k, _ := cipher.GenerateKey(32)
	jwts["HS512"] = crypto.NewJWTHS512(k)
	k, _ = cipher.GenerateKey(16)
	jwts["HS256"] = crypto.NewJWTHS256(k)
	jwts["ES256"], _ = crypto.NewJWTES256()
	jwts["ES512"], _ = crypto.NewJWTES512()
	log.Println("🔐 JWT Services: OK")

	ek = os.Getenv("FAIRHIVE_ENCRYPTION_KEY")
	if ek == "" {
		panic("encryption key is missing")
	}
	log.Println("🔑 Encryption Key: OK")

	secpath1 = os.Getenv("FAIRHIVE_API_SECURE_PATH1")
	if secpath1 == "" {
		panic("secure path #1 must be set")
	}
	secpath2 = os.Getenv("FAIRHIVE_API_SECURE_PATH2")
	if secpath2 == "" {
		panic("secure path #1 must be set")
	}
}

func NewApp() *App {
	db, err := data.NewDynamoDB(tableName, ek)
	if err != nil {
		panic(err)
	}
	return &App{
		db:       db,
		jwt:      jwts["ES256"],
		mailer:   mailer.New(os.Getenv("FAIRHIVE_GSUITE_USER"), os.Getenv("FAIRHIVE_GSUITE_PASSWORD"), "smtp.gmail.com", 587),
		wg:       sync.WaitGroup{},
		rl:       limiter.New(0.1, 10),
		secpath1: secpath1,
		secpath2: secpath2,
	}
}

var jwtregexp = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)

func (app App) register(c *gin.Context) {
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
		app.mailer.SendActivationEmail(u.Email, fmt.Sprintf("http://fairhive.io/activate/%s", token), hash)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"hash": hash,
	})
}

func (app App) activate(c *gin.Context) {
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

	c.JSON(http.StatusCreated, gin.H{
		"token":     t,
		"activated": true,
	})
}

func (app App) limit(c *gin.Context) {
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

func (app App) cors(c *gin.Context) {
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

func (app App) count(c *gin.Context) {
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

func (app App) list(c *gin.Context) {
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
			err := w.Write([]string{u.Type, u.Address, u.Email, u.UUID, fmt.Sprintf("%s", time.UnixMilli(u.Timestamp))})
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

func setupRouter(app App) *gin.Engine {
	r := gin.Default()
	t := template.Must(template.ParseFS(tfs, "templates/*"))
	r.SetHTMLTemplate(t)
	r.Use(app.cors, app.limit)
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/:path1/:path2/count", app.count)
	r.GET("/:path1/:path2/list", app.list)
	r.POST("/", app.register)
	r.POST("/activate/:token/:hash", app.activate)
	return r
}

func main() {
	app := *NewApp()
	// gin.SetMode(gin.ReleaseMode)
	r := setupRouter(app)

	var addr string
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	} else {
		addr = ":8080" // default port
	}

	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-quit
		log.Printf("🚨 Shutdown signal \"%v\" received\n", s)

		log.Printf("🚦 Here we go for a graceful Shutdown...\n")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("⚠️ HTTP server Shutdown: %v", err)
		}

		log.Printf("⏳ Waiting the end of all go-routines...")
		app.wg.Wait() // wait for all go-routines
		log.Printf("👍 go-routines are over")
		close(idleConnsClosed)
	}()

	go func() { // every 5 minutes, purge the rate limiters older than 10 minutes
		for {
			time.Sleep(5 * time.Minute)
			app.rl.Cleanup(10 * time.Minute)
		}
	}()

	log.Printf("✅ Listening and serving HTTP on %s\n", addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("👹 HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Printf("😴 Server stopped")
}
