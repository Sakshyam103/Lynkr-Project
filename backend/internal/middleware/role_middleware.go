package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "User access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func BrandOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "brand" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Brand access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserOrBrandMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "user" && role != "brand" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
