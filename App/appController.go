package Todo

import (
	gin "github.com/gin-gonic/gin"
)

// Home endpoint for Todo
func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Todo List Go Service",
	})
}
