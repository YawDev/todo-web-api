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
	"log"
	"os/exec"
	"runtime"
	auth "todo-web-api/authentication"
	app "todo-web-api/controllers"
	s "todo-web-api/storage"

	docs "todo-web-api/docs"

	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	config, err := loadConfig("config-local.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	dbConfigs := config.Database

	s.ConfigureDb(dbConfigs.UseSQLite)
	Db := s.StoreManager
	Db.Connect(dbConfigs.Username, dbConfigs.Password, dbConfigs.Host, dbConfigs.Port)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.CORSConfig.AllowedOrigins[0], config.CORSConfig.AllowedOrigins[1]},
		AllowMethods: []string{config.CORSConfig.AllowedMethods[0], config.CORSConfig.AllowedMethods[1],
			config.CORSConfig.AllowedMethods[2],
			config.CORSConfig.AllowedMethods[3]},
		AllowHeaders:     []string{config.CORSConfig.AllowedHeaders[0], config.CORSConfig.AllowedHeaders[1]},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: config.CORSConfig.AllowCredentials,
	}))

	docs.SwaggerInfo.Host = config.App.Host + ":" + config.App.Port
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
		err := openBrowser("http://" + config.App.Host + ":8080" + config.Swagger.DocPath)
		if err != nil {
			log.Println(err.Error(), err)
			panic(err)
		}
	}()
	r.Run("0.0.0.0:" + config.App.Port)
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
	r.POST("/RefreshToken", app.RefreshToken)
	r.POST("/Logout", auth.AuthMiddleware(), app.Logout)
	r.GET("/Home", app.Home)
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "darwin": // macOS
		err = exec.Command("open", url).Start() // macOS-specific command
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = exec.Command("wsl.exe", "cmd.exe", "/C", "start", url).Start()
	}
	return err
}
