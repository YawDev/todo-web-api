package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	app "todo-web-api/Controllers"
	"todo-web-api/Models"
	"todo-web-api/Storage"
	m "todo-web-api/Tests/mockmanagers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouters(userManager m.IUserMockManager) *gin.Engine {
	r := gin.Default()
	Storage.UserManager = userManager
	v1 := r.Group("/api/v1")
	{
		v1.POST("/Register", app.Register)
		r.GET("/GetUser/:id", app.GetUserById)
	}
	return r
}

func Test_Register_Url(t *testing.T) {
	router := setupRouters(&m.MockUserManager{CreateUserFn: func(user *Models.User) (int, error) {
		return 1, nil
	}})

	testUser := app.User{
		Username: "TEST_U1",
		Password: "PW1",
	}

	jsonValue, _ := json.Marshal(testUser)

	req, _ := http.NewRequest("POST", "/api/v1/Register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUser(t *testing.T) {
	router := setupRouters(&m.MockUserManager{GetUserFn: func(id int) (*Models.User, error) {
		return &Models.User{}, nil
	}})
	w := httptest.NewRecorder()

	user := Models.User{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetUser/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestExistingUser(t *testing.T) {
	router := setupRouters(&m.MockUserManager{GetUserFn: func(id int) (*Models.User, error) {
		return &Models.User{}, nil
	}})
	w := httptest.NewRecorder()

	user := Models.User{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetUser/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
