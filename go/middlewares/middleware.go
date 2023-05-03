package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHeader(header string) gin.HandlerFunc{
  return func(c *gin.Context) {
	token := c.Request.Header.Get(header)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "JWT is missing"})
		return
	}
	c.Set("jwt", token)
    c.Next()
  }
}