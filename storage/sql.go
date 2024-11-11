package storage

import (
	"fmt"
	"log"
	models "todo-web-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreDbManager struct {
}

var Context *gorm.DB

func (Db *StoreDbManager) Connect(dbUser, dbPassword, dbHost string, dbPort int) {
	var err error
	User := dbUser
	Password := dbPassword
	Port := dbPort
	Host := dbHost

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/Todo?charset=utf8mb4&parseTime=True&loc=Local", User, Password, Host, Port)
	Context, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		errMsg := "SQL Connection failed"
		log.Println(errMsg)
		panic(errMsg)
	}
	log.Println("SQL Connection Successful")
	Db.MigrateModels(Context)
}

func (Db *StoreDbManager) MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
