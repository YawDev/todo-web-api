package controllertests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	auth "todo-web-api/Authentication"
	app "todo-web-api/Controllers"
	"todo-web-api/Models"
	"todo-web-api/Storage"
	m "todo-web-api/Tests/mockmanagers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var userManager m.IUserMockManager
var listManager m.IListMockManager

func InitManagers() {
	userManager = &m.MockUserManager{
		GetUserFn: func(id int) (*Models.User, error) {
			return &Models.User{}, nil
		},
	}
	listManager = &m.MockListManager{
		GetListForUserFn: func(id int) (*Models.List, error) {
			return &Models.List{}, nil
		},
		CreateListFn: func(List *Models.List) (int, error) {
			return List.Id, nil
		}}
}

func setupListRouters(listManager m.IListMockManager, userManager m.IUserMockManager) *gin.Engine {
	r := gin.Default()
	Storage.ListManager = listManager
	Storage.UserManager = userManager
	v1 := r.Group("/api/v1")
	{
		v1.GET("/PING")
		r.POST("/CreateList/:id", auth.AuthMiddleware(), app.CreateListForUser)
		r.GET("/GetList/:userid", app.GetListByUserId)
		r.DELETE("/DeleteList/:id", auth.AuthMiddleware(), app.DeleteList)
	}
	return r
}

func TestGetListForUser(t *testing.T) {
	InitManagers()
	router := setupListRouters(listManager, userManager)
	w := httptest.NewRecorder()

	user := Models.List{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestCreateListForUser_Unauthorized(t *testing.T) {
	InitManagers()
	router := setupListRouters(listManager, userManager)
	w := httptest.NewRecorder()

	list := Models.List{}
	json, _ := json.Marshal(list)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/CreateList/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}
