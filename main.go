package main

import (
	app "todo-web-api/App"
	db "todo-web-api/Store"
	auth "todo-web-api/authentication"

	gin "github.com/gin-gonic/gin"
)

func main() {
	db.DbConnection()
	r := gin.Default()
	RouteSetup(r)
	r.Run()
}

func RouteSetup(r *gin.Engine) {
	r.GET("/Login", auth.Login)
	r.GET("/Register", auth.Register)
	r.GET("/Home", app.Home)
}
