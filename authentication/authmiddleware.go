package authentication

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		errMessage := ""

		if tokenStr == "" {
			errMessage = "access token required"

			log.Println(errMessage, errors.New(errMessage))

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMessage,
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenStr, "Bearer") {
			errMessage = "access token Bearer required"
			log.Println(errMessage, errors.New(errMessage))

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMessage,
			})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := ParseToken(tokenStr)

		if err != nil {
			log.Println(err.Error(), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}

		if _, ok := activeTokens[claims.Username]; !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token unauthorized"})
			c.Abort()
		}
		payload(claims, c)
		c.Next()
	}

}
