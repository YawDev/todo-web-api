package storage

import (
	"errors"
	"todo-web-api/messages"
	models "todo-web-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ListStore struct {
}

func (L *ListStore) CreateList(list *models.List) (ID int, err error) {
	result := Context.Create(&list)
	if result.Error != nil {
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
	}
	return list.Id, result.Error
}

func (L *ListStore) DeleteList(id int) (success bool, err error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := messages.ListNotFoundInDb
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
		return false, errors.New(errMsg)
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := "something went wrong while fetching list"
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
		return false, errors.New(errMsg)
	}

	result = Context.Delete(&list)
	if result.Error != nil {
		errMsg := "something went wrong while deleting list"
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
		return false, errors.New(errMsg)
	}
	return true, nil
}

func (L *ListStore) GetListForUser(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := messages.ListNotFoundInDb
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
		return nil, errors.New(errMsg)
	} else if result.Error != nil {
		errMsg := messages.ListNotFoundInDb
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
		return nil, errors.New(errMsg)
	}
	Context.Preload("Tasks").First(&list, id)
	return &list, nil
}

// Get List by Id
func (L *ListStore) GetList(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := messages.ListNotFoundInDb

		log.WithFields(logrus.Fields{
			"error": errMsg,
		}).Error(result.Error)
		return nil, result.Error
	} else if result.Error != nil {

		errMsg := messages.ListNotFoundInDb

		log.WithFields(logrus.Fields{
			"LoggerName": "ListStore",
			"DbContext":  "mysql",
		}).Error(result.Error.Error())
		return nil, errors.New(errMsg)
	}
	return &list, nil
}
