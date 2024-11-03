package mockmanagers

import "todo-web-api/models"

type IUserMockManager interface {
	CreateUser(user *models.User) (ID int, err error)
	DeleteUser(id int) (success bool, err error)
	GetUser(id int) (*models.User, error)
	FindExistingAccount(username string, password string) (*models.User, error)
}

type MockUserManager struct {
	CreateUserFn          func(user *models.User) (int, error)
	DeleteUserFn          func(id int) (bool, error)
	GetUserFn             func(id int) (*models.User, error)
	FindExistingAccountFn func(username string, password string) (*models.User, error)
}

func (m *MockUserManager) CreateUser(user *models.User) (int, error) {
	return m.CreateUserFn(user)
}

func (m *MockUserManager) DeleteUser(id int) (bool, error) {
	return m.DeleteUserFn(id)
}

func (m *MockUserManager) GetUser(id int) (*models.User, error) {
	return m.GetUserFn(id)
}

func (m *MockUserManager) FindExistingAccount(username string, password string) (*models.User, error) {
	return m.FindExistingAccountFn(username, password)
}
