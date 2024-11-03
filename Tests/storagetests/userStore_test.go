package storagetests

import (
	"errors"
	"testing"
	"time"
	"todo-web-api/models"
	"todo-web-api/storage"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/gorm"
)

func Test_Create_User_Successful(t *testing.T) {
	db, mock := Mock_Db_Setup()

	newUser := &models.User{
		Id:        1,
		Username:  "TestUser",
		CreatedAt: time.Now(),
		Password:  "Test_PW1",
	}

	storage.Context = db
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE Username = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(newUser.Username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(newUser.Username, newUser.Password, newUser.CreatedAt, newUser.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	success, _ := storage.UserManager.CreateUser(newUser)

	assert.Equal(t, success, 1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Create_UserExists(t *testing.T) {
	db, mock := Mock_Db_Setup()
	defer mock.ExpectClose()

	testUser := models.User{Id: 1, Username: "NewUser", CreatedAt: time.Now()}

	storage.Context = db

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE Username = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(testUser.Username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(testUser.Id, testUser.Username))

	success, err := storage.UserManager.CreateUser(&testUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	expectedError := "user exists already"
	if err != nil && err.Error() == expectedError {
		assert.Equal(t, err.Error(), expectedError)
	}
	assert.Equal(t, success, 0)
}

func Test_Get_User(t *testing.T) {
	db, mock := Mock_Db_Setup()
	defer mock.ExpectClose()

	testUser := models.User{Id: 1, Username: "NewUser"}

	storage.Context = db

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(testUser.Id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(testUser.Id, testUser.Username))

	user, _ := storage.UserManager.GetUser(testUser.Id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, user.Username, testUser.Username)
	assert.Equal(t, user.Id, testUser.Id)
}

func Test_Get_UserNotFound(t *testing.T) {
	db, mock := Mock_Db_Setup()
	defer mock.ExpectClose()

	Id := 1
	storage.Context = db

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(Id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, errMsg := storage.UserManager.GetUser(1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.EqualError(t, errors.New("user not found"), errMsg.Error())
}

func Test_Find_Existing_Account_Not_Found(t *testing.T) {
	db, mock := Mock_Db_Setup()

	newUser := &models.User{
		Id:        1,
		Username:  "TestUser",
		CreatedAt: time.Now(),
		Password:  "Test_PW1",
	}

	storage.Context = db
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE Username = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(newUser.Username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := storage.UserManager.FindExistingAccount(newUser.Username, newUser.Password)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.EqualError(t, errors.New("user not found"), err.Error())
}

func Test_Find_Existing_Account(t *testing.T) {
	db, mock := Mock_Db_Setup()

	newUser := &models.User{
		Id:        1,
		Username:  "TestUser",
		CreatedAt: time.Now(),
		Password:  "Test_PW1",
	}

	storage.Context = db
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE Username = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(newUser.Username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(newUser.Id, newUser.Username))

	user, _ := storage.UserManager.FindExistingAccount(newUser.Username, newUser.Password)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, user.Id, newUser.Id)
}
