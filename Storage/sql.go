package storage

import (
	"fmt"
	"log"
	"os"
	models "todo-web-api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreDbManager struct {
}

var Context *gorm.DB

func (Db *StoreDbManager) Connect() {
	var err error

	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	User := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASSWORD")
	Port := os.Getenv("DB_PORT")
	Host := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/Todo?charset=utf8mb4&parseTime=True&loc=Local", User, Password, Host, Port)
	Context, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
