package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	return
}

func notFound(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	return
}
