package Storage

import (
	"errors"
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/gorm"
)

type UserStore struct {
}

func (U *UserStore) CreateUser(user *models.User) (ID int, err error) {

	var existingUser models.User

	userQuery := Context.Where("Username = ?", user.Username).First(&existingUser)
	if userQuery.Error == nil {
		return 0, errors.New("user exists already")
	} else if userQuery.Error != nil && !errors.Is(userQuery.Error, gorm.ErrRecordNotFound) {
		fmt.Println("something went wrong")
		return 0, userQuery.Error
	}

	result := Context.Create(&user)
	return user.Id, result.Error
}

func (U *UserStore) DeleteUser(id int) (success bool, err error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.New("user not found")
	} else {
		fmt.Println("something went wrong")
	}

	result = Context.Delete(&user)
	if result.Error != nil {
		return false, errors.New("something went wrong while deleting User")
	}
	return true, nil
}

func (U *UserStore) GetUser(id int) (*models.User, error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if result.Error != nil {
		fmt.Println("something went wrong fetching User")
		return nil, result.Error
	}
	return &user, nil
}

func (U *UserStore) FindExistingAccount(username string, password string) (*models.User, error) {
	var user models.User
	result := Context.Where("Username = ?", username).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if result.Error != nil {
		fmt.Println("something went wrong fetching User")
		return nil, result.Error
	}
	return &user, nil
}
