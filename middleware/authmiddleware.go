package middleware

import (
	"errors"
	"net/http"
	"strings"
	auth "todo-web-api/authentication"
	l "todo-web-api/loggerutils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenStr string
		var errMessage string
		authHeader := c.GetHeader("Authorization")

		if authHeader != ""{
			l.Log.WithFields(logrus.Fields{"LoggerName": "Auth Middleware"}).Info("fetching access token from header")
			
			if(!strings.HasPrefix(authHeader, "Bearer")){
				errMessage = "access token Bearer required"
				l.ErrorLog(c.Request.Context(), http.StatusUnauthorized, errors.New(errMessage))
				c.JSON(http.StatusUnauthorized, gin.H{
				"error": errMessage,
				})
				c.Abort()
				return
			}
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}else{
			l.Log.WithFields(logrus.Fields{"LoggerName": "Auth Middleware"}).Info("fetching access token from cookies")
		    
			cookieToken, err := c.Cookie("access_token")
			if err != nil || cookieToken == ""{
				errMessage = "access token required"
				l.ErrorLog(c.Request.Context(), http.StatusUnauthorized, errors.New(errMessage))
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": errMessage,
				})
				c.Abort()
				return
			}
			tokenStr = cookieToken
		}

		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			l.ErrorLog(c.Request.Context(), http.StatusUnauthorized, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if _, ok := auth.ActiveTokens[claims.Username]; !ok {
			l.ErrorLog(c.Request.Context(), http.StatusUnauthorized, errors.New("token unauthorized"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token unauthorized"})
			c.Abort()
			return
		}

		auth.Payload(claims, c)

		l.Log.WithFields(logrus.Fields{"LoggerName": "Auth Middleware"}).Info("token authenticated")
		c.Next()
	}

}
