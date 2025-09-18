package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Server is healthy",
	})
}

func HelloWorld(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello " + name + "!",
	})
}
