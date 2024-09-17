package sqlite_db

import (
	"errors"
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) (ID int, err error) {
	result := Context.Create(&user)
	return user.Id, result.Error
}

func DeleteUser(id int) (success bool, err error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.New("user not found")
	} else {
		fmt.Println("something went wrong querying for User.")
	}

	result = Context.Delete(&user)
	if result.Error != nil {
		return false, errors.New("something went wrong while deleting User")
	}
	return true, nil
}

func GetUser(id int) (*models.User, error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if result.Error != nil {
		fmt.Println("something went wrong querying for User")
		return nil, result.Error
	}
	return &user, nil
}
