package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	h "todo-web-api/helpers"
	"todo-web-api/loggerutils"
	"todo-web-api/messages"
	models "todo-web-api/models"
	s "todo-web-api/storage"

	gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Create List endpoint for Todo
// Create List For User endpoint for Todo godoc
//
//	@BasePath		/api/v1
//	@Summary		Create List
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"id"
//	@Success		200	{object}	h.SaveResponse			"Successful"
//	@Failure		400	{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500	{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/CreateList/{id} [post]
func CreateListForUser(c *gin.Context) {
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

	result, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == "user not found, list cannot be created" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error()})
		return

	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  500,
			Message: err.Error()})
		return
	}

	existingList, err := s.ListManager.GetListForUser(id)
	if existingList != nil {
		errMsg := "there can be only 1 list per user"
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(errors.New(errMsg), err)
		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: errMsg})
		return
	} else if err != nil && !strings.Contains(err.Error(), "not found") {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error()})
		return
	}

	list := &models.List{UserId: result.Id, CreatedAt: time.Now()}
	_, err = s.ListManager.CreateList(list)
	if err != nil && err.Error() == "user not found, list cannot created" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error()})
		return

	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  500,
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.SaveResponse{
		Status:  200,
		Message: "List created successfully.",
		Id:      list.Id})
}

// Delete List For User endpoint for Todo godoc
//
//	@BasePath		/api/v1
//	@Summary		Delete List
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"id"
//	@Success		200	{object}	h.DeleteResult			"Successful"
//	@Failure		400	{object}	h.BadRequestResponse	"Bad Request"	//	Failed	due	to	bad	request	(e.g., validation error)
//	@Failure		500	{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/DeleteList/{id} [delete]
func DeleteList(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		msg := "internal error deleting list"
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, errors.New(msg))

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: msg})
		return
	}

	result, err := s.ListManager.DeleteList(id)
	if err != nil && err.Error() == "list record not found" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.ErrorResponse{
			Status:  400,
			Message: err.Error()})
		return
	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List deleted", "success": result})
}

// Fetch List By UserId godoc
//
//	@BasePath	/api/v1
//	@Summary	Get List
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int						true	"User ID"
//
//	@Success		200		{object}	h.SuccessResponse		"Successful"
//
//	@Failure		400		{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500		{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/GetList/{userid} [get]
func GetListByUserId(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("userid")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		msg := "internal error while fetching list for user"

		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, errors.New(msg))

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: msg})
	}

	user, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == "user not found" {
		msg := "user and list not found"
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, errors.New(msg))

		c.JSON(http.StatusNotFound, h.BadRequestResponse{
			Status:  404,
			Message: msg})
		return
	} else if err != nil {
		msg := "internal error while fetching list for user"
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, errors.New(msg))

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: msg})
		return
	}

	list, err := s.ListManager.GetListForUser(user.Id)
	if err != nil && err.Error() == messages.ListNotFoundInDb {
		loggerutils.ErrorLog(ctx, http.StatusNotFound,
			errors.New(messages.ListNotFoundInDb))

		c.JSON(http.StatusNotFound, h.ErrorResponse{
			Status:  404,
			Message: messages.ListNotFoundInDb})
		return
	}
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError,
			errors.New("error fetching list for user"))

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: "internal error while fetching list for user"})
		return
	}
	c.JSON(http.StatusOK, &list)
}
