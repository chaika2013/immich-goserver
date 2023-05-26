package router

import (
	"github.com/chaika2013/immich-goserver/controller"
	"github.com/gin-gonic/gin"
)

func Setup(gin *gin.Engine) {
	auth := gin.Group("/auth")
	{
		auth.POST("login", controller.Login)
		auth.POST("logout", controller.Logout)
	}

	oAuth := gin.Group("/oauth")
	{
		oAuth.POST("config", controller.GenerateConfig)
	}

	serverInfo := gin.Group("/server-info")
	serverInfo.GET("", controller.GetServerInfo)
	{
		serverInfo.GET("version", controller.GetServerVersion)
		serverInfo.GET("ping", controller.PingServer)
	}

	user := gin.Group("/user")
	{
		user.GET("count", controller.GetUserCount)

		user = user.Group("")
		user.Use(Authentication(false))
		user.GET("me", controller.GetMyUserInfo)

		user = user.Group("")
		user.Use(Authentication(true))
		user.GET("", controller.GetAllUsers)
	}

	asset := gin.Group("/asset")
	asset.Use(Authentication(false))
	{
		asset.GET(":deviceId", controller.GetUserAssetsByDeviceID)
		asset.POST("count-by-time-bucket", controller.GetAssetCountByTimeBucket)
		asset.POST("time-bucket", controller.GetAssetByTimeBucket)
		asset.GET("thumbnail/:id", controller.GetAssetThumbnail)
		asset.POST("upload", controller.UploadFile)
		asset.POST("check", controller.CheckDuplicateAsset)
	}

	jobs := gin.Group("/jobs")
	jobs.Use(Authentication(true))
	{
		jobs.GET("", controller.GetAllJobsStatus)
	}
}
