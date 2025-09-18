package routes

import (
	"github.com/chaso-pa/real-estate-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupEstateRoutes(r *gin.RouterGroup) {
	r.GET("sample", handlers.SampleCrawl)
	r.GET("crawl", handlers.CrawlSuumo)
}
