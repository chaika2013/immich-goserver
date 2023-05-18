package controller

import (
	"net/http"
	"strconv"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	req := loginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// get user
	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// verify user password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// create a session with an access token
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(http.StatusCreated, gin.H{
		"accessToken":          session.ID(),
		"userId":               strconv.FormatUint(uint64(user.ID), 10),
		"userEmail":            user.Email,
		"firstName":            user.FirstName,
		"lastName":             user.LastName,
		"profileImagePath":     "",
		"isAdmin":              user.IsAdmin,
		"shouldChangePassword": user.ShouldChangePassword,
	})
}

func Logout(c *gin.Context) {
	// delete session
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusCreated, gin.H{
		"successful":  true,
		"redirectUri": "",
	})
}
