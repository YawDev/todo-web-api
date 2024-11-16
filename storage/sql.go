package storage

import (
	"fmt"
	"todo-web-api/loggerutils"
	models "todo-web-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreDbManager struct {
}

var Context *gorm.DB

func (Db *StoreDbManager) Connect(dbUser, dbPassword, dbHost, dbPort, dbName string) {
	log := loggerutils.GetLogger()

	var err error
	User := dbUser
	Password := dbPassword
	Port := dbPort
	Host := dbHost
	Name := dbName

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", User, Password, Host, Port, dbName)
	Context, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		errMsg := "SQL Connection failed"
		log.WithFields(logrus.Fields{
			"Database": Name,
			"Host":     Host,
			"Port":     Port,
			"Error":    errMsg,
		}).Fatal(err.Error())
	}
	log.Info("SQL Connection Successful")
	Db.MigrateModels(Context)
}

func (Db *StoreDbManager) MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
