package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	app "todo-web-api/controllers"
	h "todo-web-api/helpers"
	"todo-web-api/models"
	"todo-web-api/storage"
	m "todo-web-api/tests/mockmanagers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type loginCase struct {
	existingUser *models.User
	request      h.User
	userManager  m.IUserMockManager
}

func setupRouters(userManager m.IUserMockManager) *gin.Engine {
	r := gin.Default()
	storage.UserManager = userManager
	v1 := r.Group("/api/v1")
	{
		v1.POST("/Register", app.Register)
		r.POST("/Login", app.Login)
		r.GET("/GetUser/:id", app.GetUserById)
	}
	return r
}

func Test_Login_Cases(t *testing.T) {
	existingUser := &models.User{
		Username: "u1",
	}
	existingUser.Password, _ = hashPassword("testpass1")
	var tests = []struct {
		name  string
		input loginCase
		want  int
	}{
		{
			"Successful login",
			loginCase{
				existingUser: existingUser,
				request: h.User{
					Username: "u1",
					Password: "testpass1",
				},
				userManager: &m.MockUserManager{FindExistingAccountFn: func(username, password string) (*models.User, error) {
					return existingUser, nil
				}},
			},
			200,
		},
		{
			"Failed login",
			loginCase{
				existingUser: existingUser,
				request: h.User{
					Username: "u1",
					Password: "wrongpw1",
				},
				userManager: &m.MockUserManager{FindExistingAccountFn: func(username, password string) (*models.User, error) {
					return existingUser, nil
				}},
			},
			400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := login(tt.input)
			if err != nil {
				t.Errorf(err.Error())
			}
			if ans != tt.want {
				t.Errorf("actual HTTP Status: %v, want HTTP Status: %v", ans, tt.want)
			}
			assert.Equal(t, tt.want, ans)
		})
	}
}

func login(l loginCase) (code int, errorMsg error) {
	router := setupRouters(l.userManager)
	jsonValue, _ := json.Marshal(l.request)

	req, _ := http.NewRequest("POST", "/Login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w.Code, nil
}

func Test_Register_Url(t *testing.T) {
	router := setupRouters(&m.MockUserManager{CreateUserFn: func(user *models.User) (int, error) {
		return 1, nil
	}})

	testUser := h.User{
		Username: "TEST_U1",
		Password: "PW1",
	}

	jsonValue, _ := json.Marshal(testUser)

	req, _ := http.NewRequest("POST", "/api/v1/Register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUser(t *testing.T) {
	router := setupRouters(&m.MockUserManager{GetUserFn: func(id int) (*models.User, error) {
		return &models.User{}, nil
	}})
	w := httptest.NewRecorder()

	user := models.User{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetUser/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestExistingUser(t *testing.T) {
	router := setupRouters(&m.MockUserManager{GetUserFn: func(id int) (*models.User, error) {
		return &models.User{}, nil
	}})
	w := httptest.NewRecorder()

	user := models.User{}
	json, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/GetUser/%d", 1), strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
