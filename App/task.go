package Todo

import (
	"net/http"
	"strconv"
	"time"

	models "todo-web-api/Models"
	sql "todo-web-api/sqlite_db"

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
func AddTaskToList(c *gin.Context) {
	var req SaveTask

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("listid")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := sql.GetList(id)
	if err != nil && err.Error() == "list record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	task := &models.Task{Title: req.Title, Description: req.Description, ListId: id, CreatedAt: time.Now()}

	sql.CreateTask(task, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully.",
		"Id":      result.Id,
	})
}

// Delete Task endpoint for Todo
func DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := sql.DeleteTask(id)
	if err != nil && err.Error() == "task record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully.",
		"Success": result,
	})
}

// Update Task endpoint for Todo
func UpdateTask(c *gin.Context) {
	var req SaveTask

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	task, err := sql.GetTask(id)
	if err != nil && err.Error() == "task record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	task.Title = req.Title

	if req.Description != "" {
		task.Description = req.Description
	}

	result, err := sql.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully.",
		"Id":      result,
	})
}

// Update Task endpoint for Todo
func ChangeStatus(c *gin.Context) {
	var req SetStatus

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	task, err := sql.GetTask(id)
	if err != nil && err.Error() == "task record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	task.IsCompleted = req.IsCompleted

	result, err := sql.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task status updated successfully.",
		"Id":      result,
	})
}
