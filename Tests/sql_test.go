package Tests

import (
	"todo-web-api/Storage"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Mock_Db_Setup() (*gorm.DB, sqlmock.Sqlmock) {
	dbConn, mock, err := sqlmock.New()
	dialector := mysql.New(mysql.Config{
		Conn:                      dbConn,
		DriverName:                "mysql",
		SkipInitializeWithVersion: true,
	})
	if err != nil {
		panic("Error occurred while connecting to test db: %s" + err.Error())
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("Error occurred while loading test db configs: %s" + err.Error())
	}

	CreateStores()
	return db, mock
}

func CreateStores() {
	Storage.SqlServer()

}
