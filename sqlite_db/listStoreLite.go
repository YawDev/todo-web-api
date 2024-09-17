package sqlite_db

import (
	"errors"
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/gorm"
)

func CreateList(list *models.List) (ID int, err error) {
	result := Context.Create(&list)
	return list.Id, result.Error
}

func DeleteList(id int) (success bool, err error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("List not found")
		return false, err
	} else {
		fmt.Println("Something went wrong querying for List.")
	}

	result = Context.Delete(&list)
	if result.Error != nil {
		fmt.Println("Something went wrong while deleting List.")
		return false, err
	}
	return true, nil
}

func GetListForUser(user *models.User) (*models.List, error) {
	var list models.List
	result := Context.Where("user_id = ?", user.Id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("List not found")
		return nil, result.Error
	} else if result.Error != nil {
		fmt.Println("Something went wrong querying for List.")
		return nil, result.Error
	}
	return &list, nil
}
