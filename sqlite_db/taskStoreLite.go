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
		return false, errors.New("task record not found")
	} else if result.Error != nil {
		fmt.Println("Something went wrong.")
		return false, errors.New("something went wrong")
	}

	result = Context.Delete(&task)
	if result.Error != nil {
		return false, errors.New("something went wrong")
	}
	return true, nil
}

func GetTask(id int) (*models.Task, error) {
	var task models.Task
	result := Context.First(&task, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("task record not found")
	} else if result.Error != nil {
		return nil, errors.New("something went wrong")
	}
	return &task, nil
}

func UpdateTask(task *models.Task) (ID int, err error) {
	result := Context.Save(&task)
	return task.Id, result.Error
}
