package routes

import (
	"github.com/chaso-pa/real-estate-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupStaticRoutes(r *gin.RouterGroup) {
	r.GET("/health", handlers.HealthCheck)
	r.GET("/hello", handlers.HelloWorld)
}
