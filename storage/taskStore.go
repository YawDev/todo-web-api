package storage

import (
	"errors"
	"log"
	models "todo-web-api/models"

	"gorm.io/gorm"
)

type TaskStore struct {
}

func (T *TaskStore) CreateTask(task *models.Task, listId int) (ID int, err error) {
	result := Context.Create(&task)
	return task.Id, result.Error
}

func (T *TaskStore) DeleteTask(id int) (success bool, err error) {
	var task models.Task
	result := Context.First(&task, id)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("task record not found")
		log.Println(err.Error(), err)
		return false, err
	} else if result.Error != nil {
		err := errors.New("something went wrong")
		log.Println(err.Error(), err)
		return false, err
	}

	result = Context.Delete(&task)
	if result.Error != nil {
		err := errors.New("something went wrong")
		log.Println(err.Error(), err)
		return false, err
	}
	return true, nil
}

func (T *TaskStore) GetTask(id int) (*models.Task, error) {
	var task models.Task
	result := Context.First(&task, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("task record not found")
		log.Println(err.Error(), result.Error)
		return nil, err
	} else if result.Error != nil {
		err := errors.New("something went wrong")
		log.Println(err.Error(), result.Error)
		return nil, err
	}
	return &task, nil
}

func (T *TaskStore) UpdateTask(task *models.Task) (ID int, err error) {
	result := Context.Save(&task)
	return task.Id, result.Error
}
