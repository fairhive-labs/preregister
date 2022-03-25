package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fairhive-labs/preregister/internal/crypto"
	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/fairhive-labs/preregister/internal/mailer"
	pwdgen "github.com/trendev/go-pwdgen/generator"
)

type App struct {
	db     *data.DB
	jwt    crypto.Token
	mailer mailer.Mailer
}

var jwts = map[string]crypto.Token{
	"HS512": crypto.NewJWTHS512(pwdgen.Generate(64)),
}

func init() {
	jwts["HS256"] = crypto.NewJWTHS256(pwdgen.Generate(64))
	jwts["ES256"], _ = crypto.NewJWTES256()
	jwts["ES512"], _ = crypto.NewJWTES512()
}

func NewApp(db data.DB) *App {
	return &App{&db,
		jwts["ES256"],
		mailer.NewMailer(os.Getenv("FAIRHIVE_GSUITE_USER"),
			os.Getenv("FAIRHIVE_GSUITE_PASSWORD"),
			"smtp.gmail.com",
			587),
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
	go app.mailer.SendActivationEmail(u.Email, fmt.Sprintf("http://fairhive.io/activate/%s", token), hash) //@TODO : handle graceful shutdown...
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

	err = (*app.db).Save(data.NewUser(u.Address, u.Email, u.Type))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go app.mailer.SendConfirmationEmail(u.Email) //@TODO : handle graceful shutdown...

	c.JSON(http.StatusCreated, gin.H{
		"token":     t,
		"activated": true,
	})
}

func setupRouter(app App) *gin.Engine {
	r := gin.Default()
	r.POST("/", app.register)
	r.POST("/activate/:token/:hash", app.activate)
	return r
}

func main() {
	app := *NewApp(data.MockDB) //@TODO : use dev / prod DB
	r := setupRouter(app)
	log.Fatal(r.Run())
}
