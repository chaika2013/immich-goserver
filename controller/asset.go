package controller

import (
	"bufio"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/model"
	"github.com/gin-gonic/gin"
)

type getAssetCountByTimeBucketReq struct {
	TimeGroup string `json:"timeGroup" binding:"required"`
}

func GetAssetCountByTimeBucket(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	req := getAssetCountByTimeBucketReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.TimeGroup != "month" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	timeBuckets, err := model.GetTimeBuckets(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, timeBuckets)
}

type getAssetByTimeBucketReq struct {
	TimeBucket []string `json:"timeBucket" binding:"required"`
}

func GetAssetByTimeBucket(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	req := getAssetByTimeBucketReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	assets, err := model.GetAssetsByTimeBuckets(user, req.TimeBucket)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, assets)
}

type getAssetThumbnailReq struct {
	ID string `uri:"id" binding:"required"`
}

func GetAssetThumbnail(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	_ = user

	req := getAssetThumbnailReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	format := c.Query("format")
	if format != "JPEG" && format != "WEBP" {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	// TODO: find the thumbnail

	// TODO: serve correct file
	c.File(filepath.Join(*config.ThumbnailPath, "45dddf9b-1fd3-413f-902f-798ad68c1e5b.webp"))
}

type getUserAssetsByDeviceIDReq struct {
	DeviceID string `uri:"deviceId" binding:"required"`
}

func GetUserAssetsByDeviceID(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	req := getUserAssetsByDeviceIDReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	deviceAssetIDs, err := model.GetAssetIDsByDeviceID(user, req.DeviceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, deviceAssetIDs)
}

type uploadFileReq struct {
	model.UploadFile
	AssetData     *multipart.FileHeader `form:"assetData" binding:"required"`
	LivePhotoData *multipart.FileHeader `form:"livePhotoData"`
}

func UploadFile(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	_ = user

	// TODO: optional key?
	key := c.Query("key")
	_ = key

	req := uploadFileReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// TODO: work with LivePhotoData

	// open file
	assetFile, err := req.AssetData.Open()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer assetFile.Close()

	// temp file
	tempFile, err := ioutil.TempFile(*config.UploadPath, req.AssetData.Filename+"-*")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tempFile.Close()

	// copy file
	written, err := io.Copy(bufio.NewWriter(tempFile), bufio.NewReader(assetFile))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// TODO calculate crc32

	uploadedAsset, err := model.NewAsset(user, &req.UploadFile, req.AssetData.Filename, written, 0)
	if err != nil {
		os.Remove(tempFile.Name())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, uploadedAsset)
}
