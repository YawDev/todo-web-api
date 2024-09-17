package sqlite_db

import (
	"errors"
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/gorm"
)

func CreateTask(task *models.Task, listId int) (ID int, err error) {
	result := Context.Create(&task)
	return task.Id, result.Error
}

func DeleteTask(id int) (success bool, err error) {
	var task models.Task
	result := Context.First(&task, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Task not found")
		return false, err
	} else {
		fmt.Println("Something went wrong querying for Task.")
	}

	result = Context.Delete(&task)
	if result.Error != nil {
		fmt.Println("Something went wrong while deleting Task.")
		return false, err
	}
	return true, nil
}

func GetTask(user *models.User) (*models.List, error) {
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
