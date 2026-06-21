package storagelite

import (
	"os"
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
	// SQLITE_PATH lets deployments point at a persistent volume
	// (e.g. /data/todo.db on Fly); defaults to a local file for dev.
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "todo.db"
	}
	Context, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
