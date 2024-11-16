package storagelite

import (
	l "todo-web-api/loggerutils"
	models "todo-web-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type StoreManagerLite struct {
}

var Context *gorm.DB

func (Db *StoreManagerLite) Connect(dbUser, dbPassword, dbHost, dbPort, dbName string) {
	var err error
	Context, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		errMsg := "SQLite connection failed."
		l.Log.WithFields(logrus.Fields{"LoggerName": "StoreManagerLite"}).Fatal(errMsg)
	}
	l.Log.WithFields(logrus.Fields{"LoggerName": "StoreManagerLite"}).Info("SQLite Connection Successful")
	Db.MigrateModels(Context)
}

func (Db *StoreManagerLite) MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
