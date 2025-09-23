package routes

import (
	"github.com/chaso-pa/real-estate-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUtilRoutes(r *gin.RouterGroup) {
	r.GET("/cuid", handlers.Cuid)
}
