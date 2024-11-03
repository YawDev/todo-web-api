package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	s "todo-web-api/Storage"
	models "todo-web-api/models"

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
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := s.UserManager.GetUser(id)
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

	existingList, err := s.ListManager.GetListForUser(id)
	if existingList != nil {
		errMsg := "there can be only 1 list per user"
		log.Println(errMsg, errors.New(errMsg))

		c.JSON(http.StatusBadRequest, gin.H{
			"message": errMsg,
		})
		return
	} else if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Println(err.Error(), err)

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
		log.Println(err.Error(), err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result, err := s.ListManager.DeleteList(id)
	if err != nil && err.Error() == "list record not found" {
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
		return
	} else if err != nil {
		log.Println(err.Error(), err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	list, err := s.ListManager.GetListForUser(user.Id)
	if err != nil && err.Error() == "list record not found" {
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

	c.JSON(http.StatusOK, &list)
}
