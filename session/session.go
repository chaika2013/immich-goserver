package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Setup() gin.HandlerFunc {
	return sessions.Sessions("immich_access_token", NewStore())
}
