package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"
	"time"
)

func main() {
	// Create tmp dir if not exists
	os.MkdirAll("./tmp", 0755)

	// Start pprof server
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Create CPU profile file
	f, err := os.Create("./tmp/cpu_profile.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down, stopping CPU profile...")
		pprof.StopCPUProfile()
		os.Exit(0)
	}()

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

		time.Sleep(10 * time.Second)

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
