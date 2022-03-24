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
	mailer *mailer.Mailer
}

var jwts = map[string]crypto.Token{
	"HS512": crypto.NewJWTHS512(pwdgen.Generate(64)),
}

func init() {
	jwts["HS256"] = crypto.NewJWTHS256(pwdgen.Generate(64))
	jwts["ES256"], _ = crypto.NewJWTES256()
	jwts["ES512"], _ = crypto.NewJWTES512()
}

func NewApp(db data.DB, tmplPath string) *App {
	return &App{&db,
		jwts["ES256"],
		mailer.NewMailer(os.Getenv("MAILTRAP_USER"),
			os.Getenv("MAILTRAP_PASSWORD"),
			"smtp.mailtrap.io",
			2525,
			tmplPath),
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
	token := c.Param("token")
	if !jwtregexp.MatchString(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//@TODO: get email, address, uuid from JWT
	a, e, t2 := "", "", ""
	(*app.db).Save(data.NewUser(a, e, t2))

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"activated": true,
	})
}

func setupRouter(app App) *gin.Engine {
	r := gin.Default()
	r.POST("/", app.register)
	r.POST("/activate/:token", app.activate)
	return r
}

func main() {
	app := *NewApp(data.MockDB, "internal/mailer/templates/**")
	r := setupRouter(app)
	log.Fatal(r.Run())
}
