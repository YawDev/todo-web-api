package mockmanagers

import "todo-web-api/Models"

type IUserMockManager interface {
	CreateUser(user *Models.User) (ID int, err error)
	DeleteUser(id int) (success bool, err error)
	GetUser(id int) (*Models.User, error)
	FindExistingAccount(username string, password string) (*Models.User, error)
}

type MockUserManager struct {
	CreateUserFn          func(user *Models.User) (int, error)
	DeleteUserFn          func(id int) (bool, error)
	GetUserFn             func(id int) (*Models.User, error)
	FindExistingAccountFn func(username string, password string) (*Models.User, error)
}

func (m *MockUserManager) CreateUser(user *Models.User) (int, error) {
	return m.CreateUserFn(user)
}

func (m *MockUserManager) DeleteUser(id int) (bool, error) {
	return m.DeleteUserFn(id)
}

func (m *MockUserManager) GetUser(id int) (*Models.User, error) {
	return m.GetUserFn(id)
}

func (m *MockUserManager) FindExistingAccount(username string, password string) (*Models.User, error) {
	return m.FindExistingAccountFn(username, password)
}
