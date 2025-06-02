package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Now is " + time.Now().String(),
		})
	})

	router.GET("/file", func(context *gin.Context) {
		entries, err := os.ReadDir("./tmp")
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		count := 0
		for _, entry := range entries {
			if !entry.IsDir() {
				count++
			}
		}
		context.JSON(http.StatusOK, gin.H{
			"count": count,
		})
	})

	router.POST("/", func(c *gin.Context) {
		file, err := c.FormFile("file")

		// The file cannot be received.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}

		// Retrieve file information
		extension := filepath.Ext(file.Filename)
		// Generate random file name for the new uploaded file so it doesn't override the old file with same name
		newFileName := uuid.New().String() + extension

		time.Sleep(10 * time.Second)

		// The file is received, so let's save it
		if err := c.SaveUploadedFile(file, "./tmp/"+newFileName); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}

		// File saved successfully. Return proper result
		c.JSON(http.StatusOK, gin.H{
			"message": "Your file has been successfully uploaded.",
		})
	})
	router.Run(":8080")
}
