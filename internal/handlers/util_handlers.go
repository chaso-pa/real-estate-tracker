package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucsky/cuid"
)

func Cuid(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": cuid.New(),
	})
}
