package controllers

import (
	"errors"
	http "net/http"
	"strconv"
	"strings"
	"time"
	auth "todo-web-api/authentication"
	h "todo-web-api/helpers"
	"todo-web-api/loggerutils"
	models "todo-web-api/models"

	s "todo-web-api/storage"

	msg "todo-web-api/messages"

	gin "github.com/gin-gonic/gin"
	bcr "golang.org/x/crypto/bcrypt"
)

var log = loggerutils.GetLogger()

// Login endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Login
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		h.User					true	"Login Request"
//	@Success		200		{object}	h.SuccessResponse		"Successful"
//	@Failure		400		{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500		{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/Login [post]
func Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req h.User
	var errMessage = ""
	if err := c.ShouldBindJSON(&req); err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var isLoggedIn = auth.IsTokenActive(req.Username)
	if isLoggedIn {
		errMessage = msg.AlreadyLoggedIn
		loggerutils.ErrorLog(c, http.StatusBadRequest, errors.New(errMessage))

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: errMessage})
		return
	}

	existingAccount, err := s.UserManager.FindExistingAccount(req.Username, req.Password)
	if err != nil && err.Error() == msg.AccountNotFound {

		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error()})
		return
	}

	err = bcr.CompareHashAndPassword([]byte(existingAccount.Password), []byte(req.Password))
	matchingPassword := err == nil

	if !matchingPassword {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: msg.InvalidPassword})
		return
	}

	token, err := auth.GenerateAccessToken(existingAccount.Username, existingAccount.Id)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: msg.AccessTokenError})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(existingAccount.Id, existingAccount.Username)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)
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
	loggerutils.InfoLog(ctx, http.StatusOK, msg.SuccessLogin)
	resp := h.SaveResponse{Status: 200,
		Message: "Successful Login"}
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(200,
		resp)
}

// Register endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Register
//	@Schemes
//	@Description	Create User Account
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		h.User					true	"Login Request"
//	@Success		200		{object}	h.SaveResponse			"Success"
//	@Failure		400		{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500		{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/Register [post]
func Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req h.User

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{Username: req.Username, Password: string(Hash(req.Password)), CreatedAt: time.Now()}
	id, err := s.UserManager.CreateUser(user)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: msg.SomethingWentWrong,
		})
		return
	}

	loggerutils.InfoLog(ctx, http.StatusOK, msg.SuccessUserCreate)
	c.JSON(http.StatusOK, h.SaveResponse{
		Status: 200,

		Message: msg.SuccessUserCreate,
		Id:      id,
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
//	@Success		200	{object}	h.UserResult			"Success"
//	@Success		200	{object}	h.SuccessResponse		"Successful"
//	@Failure		400	{object}	h.BadRequestResponse	"Bad Request"	//	Failed	due	to	bad	request	(e.g., validation error)
//	@Failure		500	{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/GetUser/{id} [get]
func GetUserById(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == msg.UserNotFound {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error(),
		})
		return
	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.UserResult{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

func RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	tokenStr, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not fetch refresh from cookie"})
		return
	}

	claims, err := auth.ParseRefreshToken(tokenStr)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if !auth.IsRefreshTokenActive(claims.Username) {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token unauthorized"})
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(claims.Username, claims.UserID)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: "Error while generating new access token."})
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
//	@Success		200	{object}	h.SuccessResponse		"Successful"
//	@Failure		400	{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500	{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/Logout [post]
func Logout(c *gin.Context) {
	ctx := c.Request.Context()

	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, errors.New(msg.NoTokenProvided))
		c.JSON(http.StatusUnauthorized, gin.H{"message": msg.NoTokenProvided})
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, errors.New(msg.InvalidToken))
		c.JSON(http.StatusUnauthorized, gin.H{"message": msg.InvalidToken})
	}

	auth.RemoveToken(claims.Username)
	auth.RemoveRefreshToken(claims.Username)

	loggerutils.InfoLog(ctx, http.StatusOK, msg.SuccessLogout)
	c.JSON(http.StatusOK, h.ErrorResponse{
		Status:  200,
		Message: msg.SuccessLogout,
	})
}

func Hash(password string) []byte {
	hash, err := bcr.GenerateFromPassword([]byte(password), bcr.DefaultCost)
	if err != nil {
		log.Error(err.Error())
	}
	return hash
}
