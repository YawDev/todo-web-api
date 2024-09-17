package main

import (
	app "todo-web-api/App"
	auth "todo-web-api/Authentication"
	"todo-web-api/sqlite_db"

	gin "github.com/gin-gonic/gin"
)

func main() {
	sqlite_db.Connect()
	r := gin.Default()
	RouteSetup(r)
	r.Run()
}

func RouteSetup(r *gin.Engine) {
	r.GET("/Login", auth.Login)
	r.POST("/Register", auth.Register)
	r.GET("/GetUser/:id", auth.GetUserById)
	r.GET("/Home", app.Home)
}
