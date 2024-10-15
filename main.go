// @BasePath					/api/v1
// @title						Todo Web API
// @version					1.0
// @description				Todo Web API with JWT Auth
// @host						localhost:8080
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
package main

import (
	"log"
	"os/exec"
	"runtime"
	auth "todo-web-api/Authentication"
	app "todo-web-api/Controllers"
	s "todo-web-api/Storage"

	docs "todo-web-api/docs"

	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	s.ConfigureDb()
	Db := s.StoreManager
	Db.Connect()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/")
		{
			RouteSetup(eg)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	go func() {
		err := openBrowser("http://localhost:8080/swagger/index.html")
		if err != nil {
			log.Println(err.Error(), err)
			panic(err)
		}
	}()
	r.Run(":8080")
}

func RouteSetup(r *gin.RouterGroup) {
	r.POST("/Login", app.Login)
	r.POST("/Register", app.Register)
	r.GET("/GetUser/:id", app.GetUserById)
	r.POST("/CreateList/:id", auth.AuthMiddleware(), app.CreateListForUser)
	r.DELETE("/DeleteList/:id", auth.AuthMiddleware(), app.DeleteList)
	r.GET("/GetList/:userid", app.GetListByUserId)
	r.POST("/CreateTask/:listid", auth.AuthMiddleware(), app.AddTaskToList)
	r.DELETE("/DeleteTask/:id", auth.AuthMiddleware(), app.DeleteTask)
	r.PUT("/UpdateTask/:id", auth.AuthMiddleware(), app.UpdateTask)
	r.PUT("/TaskCompleted/:id", auth.AuthMiddleware(), app.ChangeStatus)
	r.GET("/Home", auth.AuthMiddleware(), app.Home)
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "darwin": // macOS
		err = exec.Command("open", url).Start() // macOS-specific command
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	return err
}
