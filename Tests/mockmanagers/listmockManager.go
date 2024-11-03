package mockmanagers

import "todo-web-api/models"

type IListMockManager interface {
	CreateList(list *models.List) (ID int, err error)
	DeleteList(id int) (success bool, err error)
	GetListForUser(id int) (*models.List, error)
	GetList(id int) (*models.List, error)
}

type MockListManager struct {
	CreateListFn     func(list *models.List) (ID int, err error)
	DeleteListFn     func(id int) (success bool, err error)
	GetListForUserFn func(id int) (*models.List, error)
	GetListFn        func(id int) (*models.List, error)
}

func (m *MockListManager) CreateList(list *models.List) (int, error) {
	if m.CreateListFn != nil {
		return m.CreateListFn(list)
	}
	return 0, nil
}

func (m *MockListManager) DeleteList(id int) (bool, error) {
	if m.DeleteListFn != nil {
		return m.DeleteListFn(id)
	}
	return false, nil
}

func (m *MockListManager) GetListForUser(id int) (*models.List, error) {
	if m.GetListForUserFn != nil {
		return m.GetListForUserFn(id)
	}
	return nil, nil
}

func (m *MockListManager) GetList(id int) (*models.List, error) {
	if m.GetListFn != nil {
		return m.GetListFn(id)
	}
	return nil, nil
}
