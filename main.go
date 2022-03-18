package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	Address   string `json:"address" binding:"required"`
	Email     string `json:"email" binding:"required"`
	UUID      string `json:"uuid"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type" binding:"required"`
	Validated bool   `json:"validated"`
}

func main() {
	r := setupRouter()
	log.Fatal(r.Run())
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/", register)
	r.GET("/validate/:token", validate)
	return r
}

func register(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.UUID = uuid.New().String()
	u.Timestamp = time.Now().UnixMilli()
	u.Validated = false
	c.JSON(http.StatusCreated, u)
}

func validate(c *gin.Context) {
	t := c.Param("token")
	c.JSON(http.StatusOK, gin.H{
		"token":     t,
		"validated": true,
	})
}
