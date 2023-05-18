package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateConfig(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"enabled":              false,
		"passwordLoginEnabled": true,
	})
}
