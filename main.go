package main

import (
	"ats-analyzer/handlers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Create Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/analyze", handlers.AnalyzeResume)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
	}

	// Create uploads directory if it doesn't exist
	os.MkdirAll("uploads", 0755)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	logrus.Infof("Starting ATS Resume Analyzer on port %s", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
