package controller

import (
	"net/http"

	"github.com/chaika2013/immich-goserver/helper"
	"github.com/chaika2013/immich-goserver/model"
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
	user := c.MustGet("user").(*model.User)
	c.JSON(http.StatusOK, gin.H{
		"id":                   helper.StringID(user.ID),
		"email":                user.Email,
		"firstName":            user.FirstName,
		"lastName":             user.LastName,
		"createdAt":            user.CreatedAt,
		"profileImagePath":     "",
		"shouldChangePassword": user.ShouldChangePassword,
		"isAdmin":              user.IsAdmin,
		"oauthId":              "",
	})
}
