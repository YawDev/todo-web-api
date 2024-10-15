package Storage

import (
	"errors"
	"log"
	models "todo-web-api/Models"

	"gorm.io/gorm"
)

type ListStore struct {
}

func (L *ListStore) CreateList(list *models.List) (ID int, err error) {
	result := Context.Create(&list)
	log.Println(result.Error.Error(), result.Error)

	return list.Id, result.Error
}

func (L *ListStore) DeleteList(id int) (success bool, err error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := "list record not found"
		log.Println(errMsg, result.Error)

		return false, errors.New(errMsg)
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := "something went wrong"
		log.Println(errMsg, result.Error)

		return false, errors.New(errMsg)
	}

	result = Context.Delete(&list)
	if result.Error != nil {
		errMsg := "something went wrong while deleting list"
		log.Println(errMsg, result.Error)

		return false, errors.New("something went wrong")
	}
	return true, nil
}

func (L *ListStore) GetListForUser(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {

		errMsg := "list record not found"
		log.Println(errMsg, result.Error)

		return nil, errors.New(errMsg)
	} else if result.Error != nil {
		errMsg := "something went wrong"
		log.Println(errMsg, result.Error)
		return nil, result.Error
	}
	Context.Preload("Tasks").First(&list, id)
	return &list, nil
}

// Get List by Id
func (L *ListStore) GetList(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errMsg := "list record not found"
		log.Println(errMsg, result.Error)

		return nil, errors.New(errMsg)
	} else if result.Error != nil {
		errMsg := "something went wrong"
		log.Println(errMsg, result.Error)
		return nil, result.Error
	}
	return &list, nil
}
