package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	c.AbortWithStatus(http.StatusNotImplemented)
}

func validate(c *gin.Context) {
	t := c.Param("token")
	c.JSON(http.StatusOK, gin.H{
		"token":     t,
		"validated": true,
	})
}
