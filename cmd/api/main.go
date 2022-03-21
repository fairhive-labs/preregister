package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/fairhive-labs/preregister/internal/crypto"
	"github.com/fairhive-labs/preregister/internal/data"
)

type App struct {
	db  *data.DB
	jwt *crypto.JWT
}

func NewApp(db data.DB, jwt crypto.JWT) *App {
	return &App{&db, &jwt}
}

var jwtregexp = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)

func (app App) register(c *gin.Context) {
	var u data.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, u)
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
	r.GET("/activate/:token", app.activate)
	return r
}

func main() {
	app := *NewApp(data.MockDB, crypto.JWT{})
	r := setupRouter(app)
	log.Fatal(r.Run())
}
