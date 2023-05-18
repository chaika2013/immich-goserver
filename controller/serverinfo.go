package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetServerVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"major": 1,
		"minor": 55,
		"patch": 1,
	})
}

func GetServerInfo(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"diskAvailable":       "",
		"diskSize":            "",
		"diskUse":             "",
		"diskAvailableRaw":    0,
		"diskSizeRaw":         0,
		"diskUseRaw":          0,
		"diskUsagePercentage": 0,
	})
}

func PingServer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"res": "pong",
	})
}
