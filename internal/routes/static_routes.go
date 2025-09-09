package routes

import (
	"github.com/chaso-pa/gin-template/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupStaticRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck)
	r.GET("/hello", handlers.HelloWorld)
}

