package controllertests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	app "todo-web-api/controllers"
	h "todo-web-api/helpers"
	"todo-web-api/messages"
	"todo-web-api/models"
	"todo-web-api/storage"
	m "todo-web-api/tests/mockmanagers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TaskCase struct {
	id          int
	request     *h.SaveTask
	listManager m.IListMockManager
	taskManager m.ITaskMockManager
}

func setupTasksRouters(listManager m.IListMockManager, taskManager m.ITaskMockManager) *gin.Engine {
	r := gin.Default()
	storage.ListManager = listManager
	storage.TaskManager = taskManager
	v1 := r.Group("/api/v1")
	{
		v1.GET("/PING")
		r.POST("/CreateTask/:listid", app.AddTaskToList)
		r.DELETE("/DeleteTask/:id", app.DeleteTask)
		r.PUT("/UpdateTask/:id", app.UpdateTask)
		r.PUT("/TaskCompleted/:id", app.ChangeStatus)
	}
	return r
}

func Test_Task_Cases(t *testing.T) {
	task := &h.SaveTask{
		Description: "test_desc",
		Title:       "test_task",
	}
	var tests = []struct {
		name  string
		input TaskCase
		want  int
	}{
		{
			"Successful Task Add",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{CreateTaskFn: func(task *models.Task, listID int) (int, error) {
					return task.Id, nil
				}},
				listManager: &m.MockListManager{GetListFn: func(listID int) (*models.List, error) {
					return &models.List{}, nil
				}},
			},
			200,
		},
		{
			"Failed Task Add - List Not Found",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{CreateTaskFn: func(task *models.Task, listID int) (int, error) {
					return task.Id, nil
				}},
				listManager: &m.MockListManager{GetListFn: func(listID int) (*models.List, error) {
					return nil, errors.New("list record not found")
				}},
			},
			400,
		},
		{
			"Failed Task Add - Internal Error",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{CreateTaskFn: func(task *models.Task, listID int) (int, error) {
					return 0, errors.New("something went wrong")
				}},
				listManager: &m.MockListManager{GetListFn: func(listID int) (*models.List, error) {
					return &models.List{}, nil
				}},
			},
			200,
		},
		{
			"Successful Task Delete",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{DeleteTaskFn: func(id int) (bool, error) {
					return true, nil
				}},
				listManager: &m.MockListManager{},
			},
			200,
		},
		{
			"Failed Task Delete",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{DeleteTaskFn: func(id int) (bool, error) {
					return false, errors.New(messages.TaskNotFoundInDb)
				}},
				listManager: &m.MockListManager{},
			},
			400,
		},
		{
			"Successful Status Update",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{GetTaskFn: func(id int) (*models.Task, error) {
					return &models.Task{}, nil
				}},
				listManager: &m.MockListManager{},
			},
			200,
		},
		{
			"Failed Status Update",
			TaskCase{
				request: task,
				taskManager: &m.MockTaskManager{GetTaskFn: func(id int) (*models.Task, error) {
					return nil, errors.New(messages.TaskNotFoundInDb)
				}},
				listManager: &m.MockListManager{},
			},
			400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ans int
			var err error

			if strings.Contains(tt.name, "Add") {
				ans, err = CreateTask(tt.input, &h.SaveTask{
					Description: "test_desc",
					Title:       "test_task",
				})
			} else if strings.Contains(tt.name, "Delete") {
				ans, err = DeleteTask(tt.input)
			} else {
				ans, err = UpdateTask(tt.input, &h.SetStatus{
					IsCompleted: true,
				})
			}
			if err != nil {
				t.Errorf(err.Error())
			}
			if ans != tt.want {
				t.Errorf("actual HTTP Status: %v, want HTTP Status: %v", ans, tt.want)
			}

			assert.Equal(t, tt.want, ans)
		})
	}
}

func CreateTask(c TaskCase, s *h.SaveTask) (int, error) {
	router := setupTasksRouters(c.listManager, c.taskManager)
	w := httptest.NewRecorder()

	task := s
	json, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/CreateTask/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	return w.Code, nil
}

func DeleteTask(c TaskCase) (int, error) {
	router := setupTasksRouters(c.listManager, c.taskManager)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/DeleteTask/%d", 1), nil)
	router.ServeHTTP(w, req)

	return w.Code, nil
}

func UpdateTask(c TaskCase, s *h.SetStatus) (int, error) {
	router := setupTasksRouters(c.listManager, c.taskManager)
	w := httptest.NewRecorder()

	task := s
	json, _ := json.Marshal(task)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/TaskCompleted/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	return w.Code, nil
}
