package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/fairhive-labs/preregister/internal/data"
)

type App struct {
	db *data.DB
}

func NewApp(db data.DB) *App {
	return &App{&db}
}

var jwtregexp = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)

func (app App) register(c *gin.Context) {
	var u data.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Setup()
	(*app.db).Save(&u)

	c.JSON(http.StatusCreated, u)
}

func (app App) validate(c *gin.Context) {
	t := c.Param("token")
	if !jwtregexp.MatchString(t) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//@TODO: get email, address, uuid or retrieve from DB...
	u := data.User{}
	u.Validated = true
	(*app.db).Update(&u)

	c.JSON(http.StatusOK, gin.H{
		"token":     t,
		"validated": true,
	})
}

func setupRouter(app App) *gin.Engine {
	r := gin.Default()
	r.POST("/", app.register)
	r.GET("/validate/:token", app.validate)
	return r
}

func main() {
	app := *NewApp(data.MokeDB{})
	r := setupRouter(app)
	log.Fatal(r.Run())
}
