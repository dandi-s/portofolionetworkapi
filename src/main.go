package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Serve static files (dashboard)
	router.Static("/static", "./static")
	
	// Root - Serve dashboard
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "netops-integration-api",
		})
	})

	// API endpoints
	v1 := router.Group("/api/v1")
	{
		v1.POST("/agent/heartbeat", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Heartbeat received"})
		})

		v1.GET("/devices", func(c *gin.Context) {
			c.JSON(200, gin.H{"devices": []string{}})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	router.Run(":" + port)
}
