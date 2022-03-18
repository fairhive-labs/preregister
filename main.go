package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	log.Fatal(r.Run())
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/validate/:token", validate)
	return r
}

func validate(c *gin.Context) {
	c.JSON(200, gin.H{
		"validated": true,
	})
}
