package ops

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func tokenAuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token != "" {
			auth := c.GetHeader("Authorization")
			if auth != "Bearer "+token {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}
		}
		c.Next()
	}
}
