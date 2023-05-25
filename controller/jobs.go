package controller

import (
	"net/http"

	"github.com/chaika2013/immich-goserver/view"
	"github.com/gin-gonic/gin"
)

func GetAllJobsStatus(c *gin.Context) {
	// TODO check user is admin
	c.JSON(http.StatusOK, &view.AllJobs{})
}
