package router

import (
	"net/http"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return sessionBasedAuth
}

func sessionBasedAuth(c *gin.Context) {
	session := sessions.Default(c)

	// find session user
	userID := session.Get("user_id")
	if userID == nil {
		apiKeyBasedAuth(c)
		return
	}

	// check the user exist in db
	user, err := model.GetUserByID(userID.(uint))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// populate the context
	c.Set("user", user)
}

func apiKeyBasedAuth(c *gin.Context) {
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

	// populate the context
	c.Set("user", user)
}
