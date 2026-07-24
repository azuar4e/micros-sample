package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthzMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, _ := c.Get("roles")
		rolesArr := roles.([]string)

		for _, r := range rolesArr {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
