package Controllers

import (
	gin "github.com/gin-gonic/gin"
)

// Home endpoint for Todo
//
//	@BasePath	/api/v1
//	@Summary	Home
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ResponseJson	"Success"
//	@Router			/Home [get]
func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Todo List Go Service",
	})
}
