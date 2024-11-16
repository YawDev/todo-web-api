package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	h "todo-web-api/helpers"
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
	requestID, ok := ctx.Value("requestID").(string)
	if !ok {
		requestID = "unknown" // Handle missing request ID
	}

	if err != nil {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"request-ID": requestID,
			"status":     http.StatusBadRequest,
		}).Error(err.Error(), err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == "user not found" {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"request-ID": requestID,
			"status":     http.StatusBadRequest,
		}).Error(err.Error(), err)
		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error()})
		return

	} else if err != nil {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)
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
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)
		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status:  400,
			Message: err.Error()})
		return
	}

	list := &models.List{UserId: result.Id, CreatedAt: time.Now()}

	s.ListManager.CreateList(list)
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
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)
		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status:  500,
			Message: err.Error()})
		return
	}

	result, err := s.ListManager.DeleteList(id)
	if err != nil && err.Error() == "list record not found" {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)
		c.JSON(http.StatusBadRequest, h.ErrorResponse{
			Status:  400,
			Message: err.Error()})
		return
	} else if err != nil {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)
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
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	user, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == "user not found" {

		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	} else if err != nil {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)
		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
		return
	}

	list, err := s.ListManager.GetListForUser(user.Id)
	if err != nil && err.Error() == "list record not found" {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	} else if err != nil {
		log.FromContext(ctx).WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
		}).Error(err.Error(), err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &list)
}
