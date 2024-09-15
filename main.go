package main

import (
	auth "todo-web-api/auth"
	app "todo-web-api/todo"

	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	RouteSetup(r)
	r.Run()
}

func RouteSetup(r *gin.Engine) {
	r.GET("/Login", auth.Login)
	r.GET("/Register", auth.Register)
	r.GET("/Home", app.Home)
}
