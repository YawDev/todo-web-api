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
	"os/exec"
	"runtime"
	auth "todo-web-api/authentication"
	app "todo-web-api/controllers"
	"todo-web-api/loggerutils"
	"todo-web-api/middleware"
	s "todo-web-api/storage"

	docs "todo-web-api/docs"

	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	log := loggerutils.NewLogger()

	config, err := loadConfig("config-local.yaml")
	if err != nil {
		log.WithFields(logrus.Fields{
			"Error": "Unable to load config from yaml file",
		}).Fatal(err.Error())
	}
	dbConfigs := config.Database
	s.ConfigureDb(dbConfigs.UseSQLite)

	Db := s.StoreManager
	Db.Connect(dbConfigs.Username, dbConfigs.Password, dbConfigs.Host, dbConfigs.Port, dbConfigs.Name)

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
		err := openBrowser("http://" + config.App.Host + ":" + config.App.Port + config.Swagger.DocPath)
		if err != nil {
			log.WithFields(logrus.Fields{
				"Error": "Something went wrong while opening swagger on start",
			}).Error(err.Error())
		}
	}()
	r.Run("0.0.0.0:" + config.App.Port)
}

func RouteSetup(r *gin.RouterGroup) {
	r.POST("/Login", middleware.RequestIDMiddleware(), app.Login)
	r.POST("/Register", middleware.RequestIDMiddleware(), app.Register)
	r.GET("/GetUser/:id", middleware.RequestIDMiddleware(), app.GetUserById)
	r.POST("/CreateList/:id", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.CreateListForUser)
	r.DELETE("/DeleteList/:id", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.DeleteList)
	r.GET("/GetList/:userid", middleware.RequestIDMiddleware(), app.GetListByUserId)
	r.POST("/CreateTask/:listid", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.AddTaskToList)
	r.DELETE("/DeleteTask/:id", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.DeleteTask)
	r.PUT("/UpdateTask/:id", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.UpdateTask)
	r.PUT("/TaskCompleted/:id", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.ChangeStatus)
	r.POST("/RefreshToken", middleware.RequestIDMiddleware(), app.RefreshToken)
	r.POST("/Logout", middleware.RequestIDMiddleware(), auth.AuthMiddleware(), app.Logout)
	r.GET("/Home", middleware.RequestIDMiddleware(), app.Home)
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
