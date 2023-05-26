package router

import (
	"net/http"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/chaika2013/immich-goserver/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AllAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionBasedAuth(c, false)
	}
}

func AdminAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionBasedAuth(c, true)
	}
}

func sessionBasedAuth(c *gin.Context, adminOnly bool) {
	session := sessions.Default(c)

	// find session user
	userID := session.Get("user_id")
	if userID == nil {
		apiKeyBasedAuth(c, adminOnly)
		return
	}

	// check the user exist in db
	user, err := model.GetUserByID(userID.(uint))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authenticate(c, user, adminOnly)
}

func apiKeyBasedAuth(c *gin.Context, adminOnly bool) {
	apiKey := c.GetHeader("X-Api-Key")
	if apiKey == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// get user by the API key
	user, err := model.GetUserByAPIKey(apiKey)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authenticate(c, user, adminOnly)
}

func authenticate(c *gin.Context, user *model.User, adminOnly bool) {
	// check for admin
	if adminOnly && !user.IsAdmin {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// convert model User to view User
	viewUser := view.User{
		ID:                   user.ID,
		Email:                user.Email,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		CreatedAt:            user.CreatedAt,
		ShouldChangePassword: user.ShouldChangePassword,
		IsAdmin:              user.IsAdmin,
	}

	// populate the context
	c.Set("user", &viewUser)
}
