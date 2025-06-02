package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Now is " + time.Now().String(),
		})
	})

	router.POST("/", func(c *gin.Context) {
		num := rand.Intn(10000)
		c.JSON(200, gin.H{
			"message": num,
		})
	})
	router.Run(":8080")
}
