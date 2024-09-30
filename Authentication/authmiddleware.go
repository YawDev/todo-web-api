package authentication

import (
	//"log"
	//"net/http"
	//"os"
	//"time"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access token required",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenStr, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access token Bearer required",
			})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := ParseToken(tokenStr)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}

		payload(claims, c)
		c.Next()
	}

}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*User); ok && v.Username != "" {
			return true
		}
		return false
	}
}
