package sqlite_db

import (
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type StoreManagerLite struct {
}

var Context *gorm.DB

func (Db *StoreManagerLite) Initialize() {
	Db.Connect()
}

func (Db *StoreManagerLite) Connect() *gorm.DB {
	var err error
	Context, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("SQLite connection failed.")
	}
	fmt.Println("SQLite Connection Successful")
	Db.MigrateModels(Context)
	return Context
}

func (Db *StoreManagerLite) MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
