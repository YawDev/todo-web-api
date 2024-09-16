package Sqlite

import (
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Connect() {
	Db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("SQLite connection failed.")
	}
	fmt.Println("SQLite Connection Successful")
	MigrateModels(Db)
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
