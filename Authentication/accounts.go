package authentication

import (
	http "net/http"
	"strconv"
	"time"
	models "todo-web-api/Models"
	sql "todo-web-api/sqlite_db"

	bcr "golang.org/x/crypto/bcrypt"

	gin "github.com/gin-gonic/gin"
)

type User struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
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
	id, err := sql.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully.",
		"Id":      id,
	})
}

// Fetch User By Id
func GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	user, err := sql.GetUser(id)
	if err != nil && err.Error() == "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Username":  &user.Username,
		"CreatedAt": &user.CreatedAt,
	})
}

func Hash(password string) []byte {
	hash, err := bcr.GenerateFromPassword([]byte(password), bcr.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hash
}
