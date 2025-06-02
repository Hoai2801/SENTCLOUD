package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Create tmp dir if not exists
	os.MkdirAll("./tmp", 0755)

	// Gin router setup
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Now is " + time.Now().String(),
		})
	})

	router.DELETE("/file", func(context *gin.Context) {
		entries, err := os.ReadDir("./tmp")
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to read tmp directory",
			})
			return
		}

		deleted := 0
		for _, entry := range entries {
			if !entry.IsDir() {
				filePath := filepath.Join("./tmp", entry.Name())
				err := os.Remove(filePath)
				if err != nil {
					log.Printf("Failed to delete file %s: %v", filePath, err)
					continue
				}
				deleted++
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Deleted %d files", deleted),
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
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}

		extension := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + extension

		if err := c.SaveUploadedFile(file, "./tmp/"+newFileName); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your file has been successfully uploaded.",
		})
	})

	router.Run(":8080")
}
