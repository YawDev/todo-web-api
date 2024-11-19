package storagelite

import (
	"errors"
	"todo-web-api/messages"
	models "todo-web-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskStoreLite struct {
}

func (T *TaskStoreLite) CreateTask(task *models.Task, listId int) (ID int, err error) {
	result := Context.Create(&task)
	return task.Id, result.Error
}

func (T *TaskStoreLite) DeleteTask(id int) (success bool, err error) {
	var task models.Task
	result := Context.First(&task, id)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New(messages.TaskNotFoundInDb)
		log.WithFields(logrus.Fields{
			"LoggerName": "TaskStoreLite",
			"DbContext":  "sqlite",
		}).Error(err)
		return false, err
	} else if result.Error != nil {
		err := errors.New(messages.TaskQueryInternalError)
		log.WithFields(logrus.Fields{
			"LoggerName": "TaskStoreLite",
			"DbContext":  "sqlite",
		}).Error(err)
		return false, err
	}

	result = Context.Delete(&task)
	if result.Error != nil {
		err := errors.New(messages.TaskQueryInternalError)
		log.WithFields(logrus.Fields{
			"LoggerName": "TaskStoreLite",
			"DbContext":  "sqlite",
		}).Error(err)
		return false, err
	}
	return true, nil
}

func (T *TaskStoreLite) GetTask(id int) (*models.Task, error) {
	var task models.Task
	result := Context.First(&task, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New(messages.TaskNotFoundInDb)
		log.WithFields(logrus.Fields{
			"LoggerName": "TaskStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return nil, err
	} else if result.Error != nil {
		err := errors.New(messages.TaskQueryInternalError)
		log.WithFields(logrus.Fields{
			"LoggerName": "TaskStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return nil, err
	}
	return &task, nil
}

func (T *TaskStoreLite) UpdateTask(task *models.Task) (ID int, err error) {
	result := Context.Save(&task)
	return task.Id, result.Error
}
