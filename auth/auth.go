package auth

import (
	http "net/http"
	//c "golang.org/x/crypto/bcrypt"
	gin "github.com/gin-gonic/gin"
)

// Login endpoint for Todo
func Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Not Implemented Yet",
	})
}

// Register endpoint for Todo
func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Not Implemented Yet",
	})
}
