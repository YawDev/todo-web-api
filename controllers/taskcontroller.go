package controllers

import (
	"net/http"
	"strconv"
	"time"

	h "todo-web-api/helpers"
	"todo-web-api/loggerutils"
	models "todo-web-api/models"
	s "todo-web-api/storage"

	gin "github.com/gin-gonic/gin"
)

// Create Task endpoint for Todo
// Create Task By ListId godoc
//
//	@BasePath		/api/v1
//	@Summary		Create Task
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listid	path		int						true	"List ID"
//	@Param			Request	body		h.SaveTask				true	"Add Task"
//	@Success		200		{object}	h.SaveResponse		"Successful"
//	@Failure		400		{object}	h.BadRequestResponse	"Bad Request"	//	Failed	due	to	bad	request	(e.g., validation error)
//	@Failure		500		{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/CreateTask/{listid} [post]
func AddTaskToList(c *gin.Context) {
	var req h.SaveTask
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	}

	idParam := c.Param("listid")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	result, err := s.ListManager.GetList(id)
	if err != nil && err.Error() == "list record not found" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	task := &models.Task{Title: req.Title, Description: req.Description, ListId: id, CreatedAt: time.Now()}

	s.TaskManager.CreateTask(task, id)
	c.JSON(http.StatusOK, h.SaveResponse{
		Status:  200,
		Message: "Task created successfully.",
		Id:      result.Id})
}

// Delete Task For User endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Delete Task
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"id"
//	@Success		200	{object}	h.DeleteResult			"Successful"
//	@Failure		400	{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500	{object}	h.ErrorResponse			"Internal Server Error"
//	@Security		BearerAuth
//	@Router			/DeleteTask/{id} [delete]
func DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	ctx := c.Request.Context()

	id, err := strconv.Atoi(idParam)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	result, err := s.TaskManager.DeleteTask(id)
	if err != nil && err.Error() == "task record not found" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	c.JSON(http.StatusOK, h.DeleteResult{
		Status: 200,

		Message: "Task deleted successfully.",
		Success: result})
}

// Update Task For User endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Update Task
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"id"
//	@Param			Request	body		h.SaveTask				true	"Update Task"
//	@Success		200		{object}	h.SuccessResponse		"Successful"
//	@Failure		400		{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500		{object}	h.ErrorResponse			"Internal Server Error"
//	@Router			/UpdateTask/{id} [put]
func UpdateTask(c *gin.Context) {
	var req h.SaveTask
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	task, err := s.TaskManager.GetTask(id)
	if err != nil && err.Error() == "task record not found" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
		return
	}

	task.Title = req.Title

	if req.Description != "" {
		task.Description = req.Description
	}

	result, err := s.TaskManager.UpdateTask(task)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully.",
		"Id":      result,
	})
}

// Change Task Status endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Change Status Task
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		BearerAuth
//	@Param			id		path		int						true	"id"
//	@Param			Request	body		h.SetStatus				true	"Change Status"
//	@Success		200		{object}	h.SaveResponse		"Successful"
//	@Failure		400		{object}	h.BadRequestResponse	"Bad Request"
//	@Failure		500		{object}	h.ErrorResponse			"Internal Server Error"
//
//	@Router			/TaskCompleted/{id} [put]
func ChangeStatus(c *gin.Context) {
	var req h.SetStatus
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	task, err := s.TaskManager.GetTask(id)
	if err != nil && err.Error() == "task record not found" {
		loggerutils.ErrorLog(ctx, http.StatusBadRequest, err)

		c.JSON(http.StatusBadRequest, h.BadRequestResponse{
			Status: 400,

			Message: err.Error()})
		return
	} else if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
		return
	}

	task.IsCompleted = req.IsCompleted

	result, err := s.TaskManager.UpdateTask(task)
	if err != nil {
		loggerutils.ErrorLog(ctx, http.StatusInternalServerError, err)

		c.JSON(http.StatusInternalServerError, h.ErrorResponse{
			Status: 500,

			Message: err.Error()})
	}

	c.JSON(http.StatusOK, h.SaveResponse{
		Status:  http.StatusCreated,
		Message: "Task status updated successfully.",
		Id:      result,
	})
}
