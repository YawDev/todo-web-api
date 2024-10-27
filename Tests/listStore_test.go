package Tests

import (
	"testing"
	"time"
	"todo-web-api/Models"
	"todo-web-api/Storage"

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

	_, err := Storage.ListManager.CreateList(&Models.List{UserId: 1, Id: 1, CreatedAt: time.Now()})

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

// func Test_Delete_List(t *testing.T) {
// 	db, mock := Mock_Db_Setup()
// 	//defer mock.ExpectClose()
// 	Storage.Context = db

// 	listID := 1

// 	// Mock the fetch operation to ensure the list is found
// 	mock.ExpectQuery("*").
// 		WithArgs(listID).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
// 			AddRow(listID, 1, time.Now())) // Mock the row returned by First

// 	// Mock the delete operation
// 	mock.ExpectExec("DELETE FROM `lists` WHERE `lists`.`id` = \\?").
// 		WithArgs(listID).
// 		WillReturnResult(sqlmock.NewResult(0, 1)) // Indicate that 1 row was deleted

// 	_, err := Storage.ListManager.DeleteList(listID)

// 	if err != nil {
// 		t.Errorf("Failed to delete list: %s", err)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("Failed to delete list: %s", err)
// 	}
// }
