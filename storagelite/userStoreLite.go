package storagelite

import (
	"errors"
	"log"
	models "todo-web-api/models"

	"gorm.io/gorm"
)

type UserStoreLite struct {
}

func (U *UserStoreLite) CreateUser(user *models.User) (ID int, err error) {

	var existingUser models.User

	userQuery := Context.Where("Username = ?", user.Username).First(&existingUser)
	if userQuery.Error == nil {
		err := errors.New("user exists already")
		log.Println(err.Error(), err)
		return 0, err
	} else if userQuery.Error != nil && !errors.Is(userQuery.Error, gorm.ErrRecordNotFound) {
		err := errors.New("something went wrong")
		log.Println(err.Error(), err)
		return 0, userQuery.Error
	}

	result := Context.Create(&user)
	return user.Id, result.Error
}

func (U *UserStoreLite) DeleteUser(id int) (success bool, err error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("user not found")
		log.Println(err.Error(), err)
		return false, err
	} else {
		err := errors.New("something went wrong")
		log.Println(err.Error(), err)
	}

	result = Context.Delete(&user)
	if result.Error != nil {
		err := errors.New("something went wrong while deleting User")
		log.Println(err.Error(), err)
		return false, err
	}
	return true, nil
}

func (U *UserStoreLite) GetUser(id int) (*models.User, error) {
	var user models.User
	result := Context.First(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("user not found")
		log.Println(err.Error(), err)
		return nil, err
	} else if result.Error != nil {
		err := errors.New("something went wrong fetching User")
		log.Println(err.Error(), err)
		return nil, result.Error
	}
	return &user, nil
}

func (U *UserStoreLite) FindExistingAccount(username string, password string) (*models.User, error) {
	var user models.User
	result := Context.Where("Username = ?", username).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := errors.New("user not found")
		log.Println(err.Error(), err)
		return nil, err
	} else if result.Error != nil {
		err := errors.New("something went wrong fetching User")
		log.Println(err.Error(), err)
		return nil, result.Error
	}
	return &user, nil
}
