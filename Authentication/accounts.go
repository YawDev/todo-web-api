package authentication

import (
	http "net/http"
	"time"
	models "todo-web-api/Models"
	sqlite "todo-web-api/Store/Sqlite"

	bcr "golang.org/x/crypto/bcrypt"

	gin "github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Password string
}

// Login endpoint for Todo
func Login(c *gin.Context) {

	var req User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Not Implemented Yet",
	})
}

// Register endpoint for Todo
func Register(c *gin.Context) {

	var req User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{Username: req.Username, Password: string(Hash(req.Password)), CreatedAt: time.Now()}
	sqlite.CreateUser(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully.",
	})
}

func Hash(password string) []byte {
	hash, err := bcr.GenerateFromPassword([]byte(password), bcr.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hash
}
