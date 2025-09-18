package main

import (
	"log"
	"os"

	"github.com/chaso-pa/real-estate-tracker/internal/routes"
	"github.com/chaso-pa/real-estate-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()
	utils.ConDb()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	routes.SetupStaticRoutes(r.Group("/"))
	routes.SetupEstateRoutes(r.Group("/api"))

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
