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
// Root welcome page
router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "name": "Network Operations Integration API",
        "version": "1.0.0",
        "status": "running",
        "endpoints": gin.H{
            "health": "/health",
            "devices": "/api/v1/devices",
        },
    })
})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "netops-integration-api",
		})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/agent/heartbeat", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Heartbeat received"})
		})
		
		v1.GET("/devices", func(c *gin.Context) {
			c.JSON(200, gin.H{"devices": []string{}})
		})
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	router.Run(":" + port)
}
