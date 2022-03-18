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
	r.GET("/validate/:token", validate)
	r.POST("/", register)
	return r
}

func validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"validated": true,
	})
}

func register(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
