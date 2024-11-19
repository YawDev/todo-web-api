package authentication

import (
	"errors"
	"net/http"
	"strings"
	"todo-web-api/loggerutils"

	"github.com/gin-gonic/gin"
)

var log = loggerutils.GetLogger()

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		errMessage := ""

		if tokenStr == "" {
			errMessage = "access token required"
			loggerutils.ErrorLog(c.Request.Context(), http.StatusUnauthorized, errors.New(errMessage))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMessage,
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenStr, "Bearer") {
			errMessage = "access token Bearer required"
			loggerutils.ErrorLog(c.Request.Context(), http.StatusUnauthorized, errors.New(errMessage))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMessage,
			})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := ParseToken(tokenStr)

		if err != nil {
			loggerutils.ErrorLog(c.Request.Context(), http.StatusUnauthorized, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if _, ok := activeTokens[claims.Username]; !ok {
			loggerutils.ErrorLog(c.Request.Context(), http.StatusUnauthorized, errors.New("token unauthorized"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token unauthorized"})
			c.Abort()
			return
		}
		payload(claims, c)
		c.Next()
	}

}
