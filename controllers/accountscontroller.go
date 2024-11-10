package controllers

import (
	"errors"
	"fmt"
	"log"
	http "net/http"
	"strconv"
	"strings"
	"time"
	auth "todo-web-api/authentication"
	models "todo-web-api/models"

	s "todo-web-api/storage"

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
	var errMessage = ""
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var isLoggedIn = auth.IsTokenActive(req.Username)
	if isLoggedIn {
		errMessage = "User is already logged in"
		log.Println(errMessage, errors.New(errMessage))
		c.JSON(http.StatusBadRequest, gin.H{"error": errMessage})
		return
	}

	existingAccount, err := s.UserManager.FindExistingAccount(req.Username, req.Password)
	if err != nil && err.Error() == "user not found" {

		log.Println(err.Error(), err)
		c.JSON(http.StatusBadRequest, ResponseJson{Message: err.Error()})
		return
	}

	err = bcr.CompareHashAndPassword([]byte(existingAccount.Password), []byte(req.Password))
	matchingPassword := err == nil

	if !matchingPassword {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, ResponseJson{Message: "Invalid Password Credentials"})
		return
	}

	token, err := auth.GenerateAccessToken(existingAccount.Username, existingAccount.Id)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while generating access token."})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(existingAccount.Id, existingAccount.Username)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while generating refresh token."})
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

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(c.Writer, cookie)

	auth.SaveToken(existingAccount.Username, token)
	auth.SaveRefreshToken(existingAccount.Username, refreshToken)
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
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{Username: req.Username, Password: string(Hash(req.Password)), CreatedAt: time.Now()}
	id, err := s.UserManager.CreateUser(user)
	if err != nil {
		log.Println(err.Error(), err)

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
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	user, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == "user not found" {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Username":  &user.Username,
		"CreatedAt": &user.CreatedAt,
	})
}

func RefreshToken(c *gin.Context) {
	tokenStr, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not fetch refresh from cookie"})
		return
	}

	claims, err := auth.ParseRefreshToken(tokenStr)
	if err != nil {
		fmt.Println(tokenStr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(claims.Username, claims.UserID)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while generating new access token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

// Logout endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Logout
//	@Schemes
//	@Description	Logout User Account
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	ResponseJson	"Success"
//	@Router			/Logout [post]
func Logout(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided."})
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}

	auth.RemoveToken(claims.Username)
	auth.RemoveRefreshToken(claims.Username)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logout successfully.",
	})
}

func Hash(password string) []byte {
	hash, err := bcr.GenerateFromPassword([]byte(password), bcr.DefaultCost)
	if err != nil {
		log.Println(err.Error(), err)

		panic(err)
	}
	return hash
}
