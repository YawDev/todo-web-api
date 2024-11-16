package middleware

import (
	"context"
	"todo-web-api/contextkeys"
	"todo-web-api/loggerutils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func RequestIDMiddleware() gin.HandlerFunc {
	log := loggerutils.Log

	return func(c *gin.Context) {
		requestID := uuid.New()

		logger := log.Logger.WithFields(logrus.Fields{
			"RequestID": requestID,
		})
		logger.Info("Logger set in context with RequestID")

		log.Info(string(contextkeys.ContextLoggerKey) + " --middleware")
		ctx := context.WithValue(c.Request.Context(), contextkeys.ContextLoggerKey, logger)
		ctx = context.WithValue(ctx, contextkeys.ContextKeyRequestID, requestID)

		c.Request = c.Request.WithContext(ctx)

		log.WithFields(logrus.Fields{
			"ID":     requestID,
			"URL":    c.Request.Host + c.Request.RequestURI,
			"Proto":  c.Request.Proto,
			"Method": c.Request.Method,
		}).Info("Incoming Request")

		c.Next()

		log.Infof("ID: %v. Request completed.", requestID)
	}

}
