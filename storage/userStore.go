package storage

import (
	"errors"
	"todo-web-api/loggerutils"
	models "todo-web-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserStore struct {
}

var log = loggerutils.GetLogger()

func (U *UserStore) CreateUser(user *models.User) (ID int, err error) {

	var existingUser models.User

	userQuery := Context.Where("Username = ?", user.Username).First(&existingUser)
	if userQuery.Error == nil {
		err := errors.New("user exists already")
		log.WithFields(logrus.Fields{}).Error(err)
		return 0, err
	} else if userQuery.Error != nil && !errors.Is(userQuery.Error, gorm.ErrRecordNotFound) {
		err := errors.New("something went wrong creating new user")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(userQuery.Error)
		return 0, userQuery.Error
	}

	result := Context.Debug().Create(&user)
	return user.Id, result.Error
}

func (U *UserStore) DeleteUser(id int) (success bool, err error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("user not found")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
		return false, err
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("something went wrong while fetching user")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
	}

	result = Context.Delete(&user)
	if result.Error != nil {
		err := errors.New("something went wrong while deleting User")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
		return false, err
	}
	return true, nil
}

func (U *UserStore) GetUser(id int) (*models.User, error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("user not found")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
		return nil, err
	} else if result.Error != nil {
		err := errors.New("something went wrong fetching User")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (U *UserStore) FindExistingAccount(username string, password string) (*models.User, error) {
	var user models.User
	result := Context.Where("Username = ?", username).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("existing account not found")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
		return nil, err
	} else if result.Error != nil {
		err := errors.New("something went wrong fetching User")
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(result.Error)
		return nil, result.Error
	}
	return &user, nil
}
