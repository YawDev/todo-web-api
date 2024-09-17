package Todo

import (
	"net/http"
	"strconv"
	"time"

	models "todo-web-api/Models"
	sql "todo-web-api/sqlite_db"

	gin "github.com/gin-gonic/gin"
)

// Create List endpoint for Todo
func CreateListForUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := sql.GetUser(id)
	if err != nil && err.Error() == "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	existingList, err := sql.GetListForUser(id)
	if existingList != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "There can be only 1 list per user.",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	list := &models.List{UserId: result.Id, CreatedAt: time.Now()}

	sql.CreateList(list)
	c.JSON(http.StatusOK, gin.H{
		"message": "List created successfully.",
		"Id":      list.Id,
	})
}

// Delete List
func DeleteList(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := sql.DeleteList(id)
	if err != nil && err.Error() == "list record not found" {
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

	c.JSON(http.StatusOK, gin.H{"message": "List deleted", "success": result})
}

// Fetch List By UserId
func GetListByUserId(c *gin.Context) {
	idParam := c.Param("userid")
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
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	list, err := sql.GetListForUser(user.Id)
	if err != nil && err.Error() == "list record not found" {
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

	c.JSON(http.StatusOK, &list)
}
