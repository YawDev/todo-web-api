package mockmanagers

import "todo-web-api/models"

type ITaskMockManager interface {
	CreateTask(task *models.Task, listId int) (ID int, err error)
	DeleteTask(id int) (success bool, err error)
	GetTask(id int) (*models.Task, error)
	UpdateTask(task *models.Task) (ID int, err error)
}

type MockTaskManager struct {
	CreateTaskFn func(task *models.Task, listId int) (ID int, err error)
	DeleteTaskFn func(id int) (success bool, err error)
	GetTaskFn    func(id int) (*models.Task, error)
	UpdateTaskFn func(task *models.Task) (ID int, err error)
}

func (m *MockTaskManager) CreateTask(task *models.Task, listId int) (ID int, err error) {
	if m.CreateTaskFn != nil {
		return m.CreateTaskFn(task, listId)
	}
	return 0, nil
}

func (m *MockTaskManager) DeleteTask(id int) (success bool, err error) {
	if m.DeleteTaskFn != nil {
		return m.DeleteTaskFn(id)
	}
	return false, nil
}

func (m *MockTaskManager) GetTask(id int) (*models.Task, error) {
	if m.GetTaskFn != nil {
		return m.GetTaskFn(id)
	}
	return nil, nil
}

func (m *MockTaskManager) UpdateTask(task *models.Task) (int, error) {
	if m.UpdateTaskFn != nil {
		return m.UpdateTaskFn(task)
	}
	return 0, nil
}
