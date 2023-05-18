package controller

import (
	"net/http"
	"strconv"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/gin-gonic/gin"
)

func GetUserCount(c *gin.Context) {
	countAdminsOnly := c.DefaultQuery("admin", "false") == "true"
	c.JSON(http.StatusOK, gin.H{
		"userCount": model.GetUserCount(countAdminsOnly),
	})
}

func GetMyUserInfo(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	c.JSON(http.StatusOK, gin.H{
		"id":                   strconv.FormatUint(uint64(user.ID), 10),
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
