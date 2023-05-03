package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHeader() gin.HandlerFunc{
  return func(c *gin.Context) {
	token := c.Request.Header.Get("x-auth-token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "JWT is missing"})
		return
	}
	c.Set("jwt", token)
    c.Next()
  }
}