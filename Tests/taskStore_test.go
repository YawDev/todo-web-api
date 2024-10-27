package Tests

import (
	"testing"
	"time"
	"todo-web-api/Models"
	"todo-web-api/Storage"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_Create_Task(t *testing.T) {
	db, mock := Mock_Db_Setup()
	//defer mock.ExpectClose()

	Storage.Context = db

	task := Models.Task{
		Id:          1,
		Title:       "New Task",
		Description: "This is a task description",
		IsCompleted: false,
		ListId:      1,
		CreatedAt:   time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `tasks` \\(`title`,`description`,`is_completed`,`list_id`,`created_at`,`id`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
		WithArgs(task.Title, task.Description, task.IsCompleted, task.ListId, sqlmock.AnyArg(), task.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err := Storage.TaskManager.CreateTask(&task, 1)

	if err != nil {
		t.Errorf("Failed to create task: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to create task: %s", err)
	}
}

func Test_Get_Task_By_Id(t *testing.T) {
	db, mock := Mock_Db_Setup()
	Storage.Context = db

	taskID := 1
	listID := 1
	createdAt := time.Now()

	mock.ExpectQuery("SELECT \\* FROM `tasks` WHERE `tasks`.`id` = \\? ORDER BY `tasks`.`id` LIMIT \\?").
		WithArgs(taskID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "list_id", "created_at"}).
			AddRow(1, "New Task", "This is a task description", false, listID, createdAt))

	_, err := Storage.TaskManager.GetTask(1)

	if err != nil {
		t.Errorf("Failed to fetch task: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to fetch task: %s", err)
	}
}

func Test_Delete_Task(t *testing.T) {
	db, mock := Mock_Db_Setup()
	//defer mock.ExpectClose()

	Storage.Context = db

	taskID := 1
	listID := 1
	createdAt := time.Now()

	mock.ExpectQuery("SELECT \\* FROM `tasks` WHERE `tasks`.`id` = \\? ORDER BY `tasks`.`id` LIMIT \\?").
		WithArgs(taskID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "list_id", "created_at"}).
			AddRow(1, "New Task", "This is a task description", false, listID, createdAt))

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `tasks` WHERE `tasks`.`id` = \\?").
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	success, err := Storage.TaskManager.DeleteTask(1)

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
