package controllers

import (
	"errors"
	http "net/http"
	"strconv"
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

		loggerutils.ErrorLog(ctx, http.StatusNotFound, err)

		c.JSON(http.StatusNotFound, h.NotFoundResponse{
			Status:  404,
			Message: msg.AccountNotFound})
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

	http.SetCookie(c.Writer, &http.Cookie{
	Name:     "access_token",
	Value:    token,
	Path:     "/",
	Domain:   "",
	MaxAge:   3600,
	HttpOnly: true,
	Secure:   true,
	SameSite: http.SameSiteNoneMode,
})

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
	resp := h.AuthStatusResponse{Status: 200,
		Message: "Successful Login",
		User:  h.UserContext{
			Username: existingAccount.Username,
			Id: existingAccount.Id,
		},
	}
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
		loggerutils.ErrorLog(ctx, http.StatusNotFound, err)
		c.JSON(http.StatusNotFound, h.NotFoundResponse{
			Status:  404,
			Message: msg.UserNotFound,
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
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not fetch refresh from cookie"})
		return
	}

	claims, err := auth.ParseRefreshToken(tokenStr)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg.InvalidRefreshToken})
		return
	}

	if !auth.IsRefreshTokenActive(claims.Username) {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg.UnauthorizedRefreshToken})
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

	cookieToken, err := c.Cookie("access_token")
	if err != nil  || cookieToken == ""{
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, errors.New(msg.NoTokenProvided))
		c.JSON(http.StatusUnauthorized, gin.H{"message": msg.NoTokenProvided})
	}

	claims, err := auth.ParseToken(cookieToken)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, errors.New(msg.InvalidToken))
		c.JSON(http.StatusUnauthorized, gin.H{"message": msg.InvalidToken})
	}

	auth.RemoveToken(claims.Username)
	auth.RemoveRefreshToken(claims.Username)

	//Remove tokens from browser cookies
	c.SetCookie("access_token", "", -1, "/", "", true, true) 
	c.SetCookie("refresh_token", "", -1, "/", "", true, true) 

	loggerutils.InfoLog(ctx, http.StatusOK, msg.SuccessLogout)
	c.JSON(http.StatusOK, h.ErrorResponse{
		Status:  200,
		Message: msg.SuccessLogout,
	})
}

func AuthStatus(c *gin.Context) {
	ctx := c.Request.Context()
	tokenStr, err := c.Cookie("access_token")
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not fetch token from cookie"})
		return
	}

	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !auth.IsTokenActive(claims.Username) {
		loggerutils.InfoLog(ctx, http.StatusUnauthorized, msg.InvalidToken)
		c.JSON(http.StatusUnauthorized, gin.H{"Message": msg.InvalidToken, "Status":401})
		return
	}

	message := "User is currently logged in"
	loggerutils.InfoLog(c, http.StatusOK, message)

	c.JSON(http.StatusOK, h.AuthStatusResponse{
			Status:  200,
			Message: message,
			User: h.UserContext{
				Username: claims.Username,
				Id:  claims.UserID,

			}})
}

func Hash(password string) []byte {
	hash, err := bcr.GenerateFromPassword([]byte(password), bcr.DefaultCost)
	if err != nil {
		log.Error(err.Error())
	}
	return hash
}
