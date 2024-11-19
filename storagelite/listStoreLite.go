package storagelite

import (
	"errors"
	"todo-web-api/loggerutils"
	"todo-web-api/messages"
	models "todo-web-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var log = loggerutils.GetLogger()

type ListStoreLite struct {
}

func (L *ListStoreLite) CreateList(list *models.List) (ID int, err error) {
	result := Context.Create(&list)
	log.WithFields(logrus.Fields{
		"LoggerName": "ListStoreLite",
		"DbContext":  "sqlite",
	}).Error(result.Error)

	return list.Id, result.Error
}

func (L *ListStoreLite) DeleteList(id int) (success bool, err error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := messages.ListNotFoundInDb
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return false, errors.New(errMsg)
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := messages.ListQueryInternalError
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return false, errors.New(errMsg)
	}

	result = Context.Delete(&list)
	if result.Error != nil {
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return false, errors.New(messages.ListQueryInternalError)
	}
	return true, nil
}

func (L *ListStoreLite) GetListForUser(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {

		errMsg := messages.ListNotFoundInDb
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return nil, errors.New(errMsg)
	} else if result.Error != nil {
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return nil, errors.New(messages.ListQueryInternalError)
	}
	Context.Preload("Tasks").First(&list, id)
	return &list, nil
}

// Get List by Id
func (L *ListStoreLite) GetList(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := messages.ListNotFoundInDb
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return nil, errors.New(errMsg)
	} else if result.Error != nil {
		errMsg := messages.ListQueryInternalError
		log.WithFields(logrus.Fields{
			"LoggerName": "ListStoreLite",
			"DbContext":  "sqlite",
		}).Error(result.Error)
		return nil, errors.New(errMsg)
	}
	return &list, nil
}
