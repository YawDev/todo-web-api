package server

import (
	"os/exec"
	"runtime"
	app "todo-web-api/controllers"
	"todo-web-api/loggerutils"
	"todo-web-api/middleware"
	store "todo-web-api/storage"

	docs "todo-web-api/docs"

	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	config *Config
	logger *loggerutils.LogUtil
}

func NewService(c *Config, l *loggerutils.LogUtil) *Service {
	return &Service{
		config: c,
		logger: l,
	}
}

func (s *Service) Start(r *gin.Engine) {
	s.connectToSQL()
	s.corsConfiguration(r)
	if s.config.Swagger.Enabled {
		s.swaggerSetup(r)
	}
	if err := r.Run(s.config.App.Host + ":" + s.config.App.Port); err != nil {
		s.logger.WithFields(logrus.Fields{"Error": "Unable to start application",
			"Port": s.config.App.Port,
		}).Fatal(err)
	}
}

func (s *Service) connectToSQL() {
	dbConfigs := s.config.Database
	store.ConfigureDb(dbConfigs.UseSQLite)
	Db := store.StoreManager
	Db.Connect(dbConfigs.Username, dbConfigs.Password, dbConfigs.Host, dbConfigs.Port, dbConfigs.Name)
}

func (s *Service) corsConfiguration(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.CORSConfig.AllowedOrigins,
		AllowMethods:     s.config.CORSConfig.AllowedMethods,
		AllowHeaders:     s.config.CORSConfig.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: s.config.CORSConfig.AllowCredentials,
	}))
}

func (s *Service) swaggerSetup(r *gin.Engine) {
	docs.SwaggerInfo.Host = s.config.App.Host + ":" + s.config.App.Port
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/")
		{
			s.RouteSetup(eg)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	go func() {
		err := s.openBrowser("http://" + s.config.App.Host + ":" + s.config.App.Port + s.config.Swagger.DocPath)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"Error": "Something went wrong while opening swagger on start",
			}).Error(err.Error())
		}
	}()
}

func (s *Service) RouteSetup(r *gin.RouterGroup) {

	v1 := r.Group("/", middleware.RequestIDMiddleware())
	{
		v1.GET("/Home", middleware.RequestIDMiddleware(), app.Home)
		v1.GET("/AuthStatus", middleware.RequestIDMiddleware(), app.AuthStatus)
		v1.POST("/Login", middleware.RequestIDMiddleware(), app.Login)
		v1.POST("/Register", middleware.RequestIDMiddleware(), app.Register)
		v1.POST("/RefreshToken", middleware.RequestIDMiddleware(), app.RefreshToken)
	}

	auth := r.Group("/", middleware.RequestIDMiddleware(), middleware.AuthMiddleware())
	{
		auth.GET("/GetUser/:id", app.GetUserById)
		auth.POST("/CreateList/:id", app.CreateListForUser)
		auth.GET("/GetList/:userid", app.GetListByUserId)
		auth.DELETE("/DeleteList/:id", app.DeleteList)
		auth.POST("/CreateTask/:listid", app.AddTaskToList)
		auth.DELETE("/DeleteTask/:id", app.DeleteTask)
		auth.PUT("/UpdateTask/:id", app.UpdateTask)
		auth.PUT("/TaskCompleted/:id", app.ChangeStatus)
		auth.POST("/Logout", app.Logout)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (s *Service) openBrowser(url string) error {
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
