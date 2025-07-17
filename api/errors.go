package api

import (
	"github.com/gin-gonic/gin"
)

func errResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
}
