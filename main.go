package main

import (
	auth "todo-web-api/auth"

	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/Login", auth.Login)
	r.GET("/Register", auth.Register)
	r.Run()
}
