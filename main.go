// @BasePath					/api/v1
// @title						Todo.Service
// @version					1.0
// @description				Todo.Service
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
package main

import (
	"todo-web-api/loggerutils"
	s "todo-web-api/server"

	gin "github.com/gin-gonic/gin"
)

var Log *loggerutils.LogUtil = loggerutils.NewLogger()

func main() {
	r := gin.Default()
	s := &s.Service{}
	s.Start(r)

}
