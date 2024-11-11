package storagelite

import (
	"log"
	models "todo-web-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type StoreManagerLite struct {
}

var Context *gorm.DB

func (Db *StoreManagerLite) Connect(dbUser, dbPassword, dbHost, dbPort string) {
	var err error
	Context, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		errMsg := "SQLite connection failed."
		log.Println(errMsg)
		panic(errMsg)
	}
	log.Println("SQLite Connection Successful")
	Db.MigrateModels(Context)
}

func (Db *StoreManagerLite) MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
}
