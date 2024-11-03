package storagetests

import (
	"testing"
	"time"
	"todo-web-api/Storage"
	"todo-web-api/models"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_Create_List(t *testing.T) {
	db, mock := Mock_Db_Setup()
	//defer mock.ExpectClose()

	Storage.Context = db

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `lists` \\(`user_id`,`created_at`,`id`\\) VALUES \\(\\?,\\?,\\?\\)").
		WithArgs(1, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err := Storage.ListManager.CreateList(&models.List{UserId: 1, Id: 1, CreatedAt: time.Now()})

	if err != nil {
		t.Errorf("Failed to create list: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to create list: %s", err)
	}
}

func Test_Get_List_By_Id(t *testing.T) {
	db, mock := Mock_Db_Setup()
	//defer mock.ExpectClose()

	listID := 1
	userID := 1
	createdAt := time.Now()

	Storage.Context = db

	mock.ExpectQuery("SELECT \\* FROM `lists` WHERE `lists`.`id` = \\? ORDER BY `lists`.`id` LIMIT \\?").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
			AddRow(listID, userID, createdAt))

	_, err := Storage.ListManager.GetList(1)

	if err != nil {
		t.Errorf("Failed to fetch list: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to fetch list: %s", err)
	}
}

func Test_Delete_List(t *testing.T) {
	db, mock := Mock_Db_Setup()
	//defer mock.ExpectClose()

	Storage.Context = db

	listID := 1
	userID := 1
	createdAt := time.Now()

	mock.ExpectQuery("SELECT \\* FROM `lists` WHERE `lists`.`id` = \\? ORDER BY `lists`.`id` LIMIT \\?").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
			AddRow(listID, userID, createdAt))

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `lists` WHERE `lists`.`id` = \\?").
		WithArgs(listID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	success, err := Storage.ListManager.DeleteList(1)

	if err != nil {
		t.Errorf("Failed to fetch list: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to fetch list: %s", err)
	}

	if success {
		t.Log("Delete successful")
	}

	assert.True(t, success)
}
