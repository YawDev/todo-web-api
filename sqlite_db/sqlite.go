package sqlite_db

import (
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Context *gorm.DB

func Connect() *gorm.DB {
	var err error
	Context, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("SQLite connection failed.")
	}
	fmt.Println("SQLite Connection Successful")
	MigrateModels(Context)
	return Context
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
