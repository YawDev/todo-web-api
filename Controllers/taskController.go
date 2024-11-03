package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	s "todo-web-api/Storage"
	models "todo-web-api/models"

	gin "github.com/gin-gonic/gin"
)

type SaveTask struct {
	Title       string `binding:"required"`
	Description string
}

type SetStatus struct {
	IsCompleted bool `binding:"required"`
}

// Create Task endpoint for Todo
// Create Task By ListId godoc
//
//	@BasePath		/api/v1
//	@Summary		Create Task
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listid	path		int				true	"List ID"
//	@Param			Request	body		SaveTask		true	"Add Task"
//	@Success		200		{object}	ResponseJson	"Success"
//	@Router			/CreateTask/{listid} [post]
func AddTaskToList(c *gin.Context) {
	var req SaveTask

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("listid")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := s.ListManager.GetList(id)
	if err != nil && err.Error() == "list record not found" {
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

	task := &models.Task{Title: req.Title, Description: req.Description, ListId: id, CreatedAt: time.Now()}

	s.TaskManager.CreateTask(task, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully.",
		"Id":      result.Id,
	})
}

// Delete Task For User endpoint for Todo godoc
//
//	@BasePath	/api/v1
//	@Summary	Delete Task
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"id"
//	@Success		200	{object}	ResponseJson	"Success"
//	@Security		BearerAuth
//	@Router			/DeleteTask/{id} [delete]
func DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := s.TaskManager.DeleteTask(id)
	if err != nil && err.Error() == "task record not found" {
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
		"message": "Task deleted successfully.",
		"Success": result,
	})
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
//	@Param			id		path		int				true	"id"
//	@Param			Request	body		SaveTask		true	"Update Task"
//	@Success		200		{object}	ResponseJson	"Success"
//	@Router			/UpdateTask/{id} [put]
func UpdateTask(c *gin.Context) {
	var req SaveTask

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	task, err := s.TaskManager.GetTask(id)
	if err != nil && err.Error() == "task record not found" {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	task.Title = req.Title

	if req.Description != "" {
		task.Description = req.Description
	}

	result, err := s.TaskManager.UpdateTask(task)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
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
//	@Param			id		path		int				true	"id"
//	@Param			Request	body		SetStatus		true	"Change Status"
//	@Success		200		{object}	ResponseJson	"Success"
//	@Router			/TaskCompleted/{id} [put]
func ChangeStatus(c *gin.Context) {
	var req SetStatus

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	task, err := s.TaskManager.GetTask(id)
	if err != nil && err.Error() == "task record not found" {
		log.Println(err.Error(), err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	task.IsCompleted = req.IsCompleted

	result, err := s.TaskManager.UpdateTask(task)
	if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task status updated successfully.",
		"Id":      result,
	})
}
