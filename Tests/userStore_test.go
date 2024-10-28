package Tests

import (
	"testing"
	"time"
	"todo-web-api/Models"
	"todo-web-api/Storage"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_Create_UserExists(t *testing.T) {
	db, mock := Mock_Db_Setup()
	//defer mock.ExpectClose()

	testUser := Models.User{Username: "NewUser", CreatedAt: time.Now()}

	Storage.Context = db

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE Username = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(testUser.Username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "TestUser"))

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(testUser.Username, testUser.CreatedAt).
		WillReturnResult(sqlmock.NewResult(2, 1)) // Returning new user ID
	mock.ExpectCommit()
}
