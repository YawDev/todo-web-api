package authentication

import (
	http "net/http"
	"strconv"
	"time"
	models "todo-web-api/Models"
	s "todo-web-api/Storage"

	gin "github.com/gin-gonic/gin"
	bcr "golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type ResponseJson struct {
	Message string `json:"message" example:"Success"`
}

// Login endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Login
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		User			true	"Login Request"
//	@Success		200		{object}	ResponseJson	"Successful"
//	@Router			/Login [post]
func Login(c *gin.Context) {

	var req User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingAccount, err := s.UserManager.FindExistingAccount(req.Username, req.Password)
	if err != nil && err.Error() == "user not found" {
		c.JSON(http.StatusBadRequest, ResponseJson{Message: err.Error()})
		return
	}

	err = bcr.CompareHashAndPassword([]byte(existingAccount.Password), []byte(req.Password))
	matchingPassword := err == nil

	if !matchingPassword {
		c.JSON(http.StatusBadRequest, ResponseJson{Message: "Invalid Password Credentials"})
		return
	}

	token, err := GenerateAccessToken(existingAccount.Username, existingAccount.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while generating access token."})
		return
	}

	c.SetCookie(
		"access_token",
		token,
		3600,
		c.Request.RequestURI,
		"localhost",
		true,
		true,
	)

	resp := ResponseJson{Message: "Successful Login"}
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(200, resp)
}

// Register endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Register
//	@Schemes
//	@Description	Create User Account
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		User			true	"Login Request"
//	@Success		200		{object}	ResponseJson	"Success"
//	@Router			/Register [post]
func Register(c *gin.Context) {

	var req User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{Username: req.Username, Password: string(Hash(req.Password)), CreatedAt: time.Now()}
	id, err := s.UserManager.CreateUser(user)
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
//
//	@BasePath	/api/v1
//	@Summary	GetUserById
//	@Schemes
//	@Description	Fetch User Account
//	@Param			id	path	int	true	"id"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ResponseJson	"Success"
//	@Router			/GetUser/{id} [get]
func GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	user, err := s.UserManager.GetUser(id)
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
