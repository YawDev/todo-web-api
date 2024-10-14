package Storage

import (
	"errors"
	"fmt"
	models "todo-web-api/Models"

	"gorm.io/gorm"
)

type ListStore struct {
}

func (L *ListStore) CreateList(list *models.List) (ID int, err error) {
	result := Context.Create(&list)
	return list.Id, result.Error
}

func (L *ListStore) DeleteList(id int) (success bool, err error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.New("list record not found")
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.New("something went wrong")
	}

	result = Context.Delete(&list)
	if result.Error != nil {
		fmt.Println("something went wrong while deleting list")
		return false, errors.New("something went wrong")
	}
	return true, nil
}

func (L *ListStore) GetListForUser(id int) (*models.List, error) {
	var list models.List
	result := Context.First(&list, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("list record not found")
	} else if result.Error != nil {
		fmt.Println("something went wrong")
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
		return nil, errors.New("list record not found")
	} else if result.Error != nil {
		fmt.Println("something went wrong.")
		return nil, result.Error
	}
	return &list, nil
}