package controllertests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	app "todo-web-api/controllers"
	"todo-web-api/models"
	"todo-web-api/storage"
	m "todo-web-api/tests/mockmanagers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var userManager m.IUserMockManager
var listManager m.IListMockManager

func InitManagersDefault() {
	userManager = &m.MockUserManager{
		GetUserFn: func(id int) (*models.User, error) {
			return &models.User{}, nil
		},
	}
	listManager = &m.MockListManager{
		GetListForUserFn: func(id int) (*models.List, error) {
			return &models.List{}, nil
		},
		CreateListFn: func(List *models.List) (int, error) {
			return List.Id, nil
		}}
}

func setupListRouters(listManager m.IListMockManager, userManager m.IUserMockManager) *gin.Engine {
	r := gin.Default()
	storage.ListManager = listManager
	storage.UserManager = userManager
	v1 := r.Group("/api/v1")
	{
		v1.GET("/PING")
		r.POST("/CreateList/:id", app.CreateListForUser)
		r.GET("/GetList/:userid", app.GetListByUserId)
		r.DELETE("/DeleteList/:id", app.DeleteList)
	}
	return r
}

func TestGetListForUser(t *testing.T) {
	InitManagersDefault()
	router := setupListRouters(listManager, userManager)
	w := httptest.NewRecorder()

	user := models.List{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGetListForUser_NotFound(t *testing.T) {
	InitManagersDefault()
	router := setupListRouters(listManager, userManager)
	w := httptest.NewRecorder()

	user := models.List{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestCreateListForUser(t *testing.T) {
	InitManagersDefault()
	router := setupListRouters(&m.MockListManager{
		GetListForUserFn: func(id int) (*models.List, error) {
			return nil, nil
		}}, userManager)
	w := httptest.NewRecorder()

	list := models.List{}
	json, _ := json.Marshal(list)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/CreateList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestCreateListForUser_NotFound(t *testing.T) {
	InitManagersDefault()
	router := setupListRouters(listManager, &m.MockUserManager{
		GetUserFn: func(id int) (*models.User, error) {
			return nil, errors.New("user not found")
		}})
	w := httptest.NewRecorder()

	user := models.List{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/CreateList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCreateList_AlreadyExists(t *testing.T) {
	InitManagersDefault()
	router := setupListRouters(&m.MockListManager{
		GetListForUserFn: func(id int) (*models.List, error) {
			return &models.List{}, nil
		}}, userManager)
	w := httptest.NewRecorder()

	list := models.List{}
	json, _ := json.Marshal(list)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/CreateList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestDeleteList(t *testing.T) {
	router := setupListRouters(listManager, userManager)
	w := httptest.NewRecorder()

	list := models.List{}
	json, _ := json.Marshal(list)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/DeleteList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteList_ListNotFound(t *testing.T) {
	router := setupListRouters(&m.MockListManager{
		DeleteListFn: func(id int) (bool, error) {
			return false, errors.New("list record not found")
		}}, userManager)

	w := httptest.NewRecorder()

	list := models.List{}
	json, _ := json.Marshal(list)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/DeleteList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
