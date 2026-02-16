package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	"portofolionetworkapi/internal/database"
	"portofolionetworkapi/internal/handlers"
)

func main() {
	godotenv.Load()

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Println("Warning: Migration failed:", err)
	}

	router := gin.Default()

	// CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Static files
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "netops-integration-api"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/agent/heartbeat", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Heartbeat received"})
		})

		devices := v1.Group("/devices")
		{
			devices.GET("", handlers.ListDevices)
			devices.POST("", handlers.CreateDevice)
			devices.PUT("/:id", handlers.UpdateDevice)
			devices.DELETE("/:id", handlers.DeleteDevice)
		}
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
