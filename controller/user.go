package controller

import (
	"net/http"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/chaika2013/immich-goserver/view"
	"github.com/gin-gonic/gin"
)

func GetUserCount(c *gin.Context) {
	countAdminsOnly := c.DefaultQuery("admin", "false") == "true"
	userCount, err := model.GetUserCount(countAdminsOnly)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userCount": userCount,
	})
}

func GetMyUserInfo(c *gin.Context) {
	user := c.MustGet("user").(*view.User)
	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	user := c.MustGet("user").(*view.User)
	isAll := c.DefaultQuery("isAll", "false") == "true"

	users, err := model.GetAllUsers(user.ID, isAll)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
