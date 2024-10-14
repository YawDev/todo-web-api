package Todo

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	models "todo-web-api/Models"
	s "todo-web-api/Storage"

	gin "github.com/gin-gonic/gin"
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
//	@Param			id	path		int				true	"id"
//	@Success		200	{object}	ResponseJson	"Success"
//	@Router			/CreateList/{id} [post]
func CreateListForUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := s.UserManager.GetUser(id)
	if err != nil && err.Error() == "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	existingList, err := s.ListManager.GetListForUser(id)
	if existingList != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "There can be only 1 list per user.",
		})
		return
	} else if err != nil && !strings.Contains(err.Error(), "not found") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	list := &models.List{UserId: result.Id, CreatedAt: time.Now()}

	s.ListManager.CreateList(list)
	c.JSON(http.StatusOK, gin.H{
		"message": "List created successfully.",
		"Id":      list.Id,
	})
}

// Delete List For User endpoint for Todo godoc
//
//	@BasePath		/api/v1
//	@Summary		Delete List
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int				true	"id"
//	@Success		200	{object}	ResponseJson	"Success"
//	@Router			/DeleteList/{id} [delete]
func DeleteList(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := s.ListManager.DeleteList(id)
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

// Fetch List By UserId godoc
//
//	@BasePath	/api/v1
//	@Summary	Get List
//	@Schemes
//	@Description	Sign-In with user credentials, for generated access token
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int				true	"User ID"
//	@Success		200		{object}	ResponseJson	"Success"
//	@Router			/GetList/{userid} [get]
func GetListByUserId(c *gin.Context) {
	idParam := c.Param("userid")
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
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	list, err := s.ListManager.GetListForUser(user.Id)
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
