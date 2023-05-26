package controller

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/helper"
	"github.com/chaika2013/immich-goserver/model"
	"github.com/chaika2013/immich-goserver/pipeline"
	"github.com/chaika2013/immich-goserver/view"
	"github.com/gin-gonic/gin"
)

type getAssetCountByTimeBucketReq struct {
	TimeGroup string `json:"timeGroup" binding:"required"`
}

func GetAssetCountByTimeBucket(c *gin.Context) {
	user := c.MustGet("user").(*view.User)

	req := getAssetCountByTimeBucketReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.TimeGroup != "month" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	timeBuckets, err := model.GetTimeBuckets(user.ID)
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
	user := c.MustGet("user").(*view.User)

	req := getAssetByTimeBucketReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	assets, err := model.GetAssetsByTimeBuckets(user.ID, req.TimeBucket)
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
	user := c.MustGet("user").(*view.User)

	req := getAssetThumbnailReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	format := c.Query("format")
	if format != "JPEG" && format != "WEBP" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.File(filepath.Join(*config.ThumbnailPath, helper.StringID(user.ID),
		fmt.Sprintf("%s.%s", req.ID, strings.ToLower(format))))
}

type getUserAssetsByDeviceIDReq struct {
	DeviceID string `uri:"deviceId" binding:"required"`
}

func GetUserAssetsByDeviceID(c *gin.Context) {
	user := c.MustGet("user").(*view.User)

	req := getUserAssetsByDeviceIDReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	deviceAssetIDs, err := model.GetAssetIDsByDeviceID(user.ID, req.DeviceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, deviceAssetIDs)
}

type uploadFileReq struct {
	view.UploadFile
	AssetData     *multipart.FileHeader `form:"assetData" binding:"required"`
	LivePhotoData *multipart.FileHeader `form:"livePhotoData"`
}

func UploadFile(c *gin.Context) {
	user := c.MustGet("user").(*view.User)

	// TODO: optional key?
	key := c.Query("key")
	if key != "" {
		c.AbortWithError(http.StatusNotImplemented, errors.New("implement key"))
		return
	}

	req := uploadFileReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// TODO: work with LivePhotoData
	if req.LivePhotoData != nil {
		c.AbortWithError(http.StatusNotImplemented, errors.New("implement LivePhotoData"))
		return
	}

	// open file
	assetFile, err := req.AssetData.Open()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer assetFile.Close()

	// temp file
	uploadPath, err := helper.MakeUserUploadDir(user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tempFile, err := os.CreateTemp(uploadPath, req.AssetData.Filename+"-*")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, fileName := filepath.Split(tempFile.Name())
	defer tempFile.Close()

	// copy file and calculate crc32
	crc32Writer := helper.NewCRC32Writer(crc32.Castagnoli, bufio.NewWriter(tempFile))
	written, err := io.Copy(crc32Writer, bufio.NewReader(assetFile))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := crc32Writer.Flush(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// create new asset
	asset, err := model.NewUploadAsset(user.ID, &req.UploadFile, req.AssetData.Filename,
		written, crc32Writer.Sum(), fileName)
	if err != nil {
		os.Remove(tempFile.Name())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// start processing file
	pipeline.Enqueue(asset.ID, pipeline.AllJobs)

	c.JSON(http.StatusCreated, gin.H{
		"id":        helper.StringID(asset.ID),
		"duplicate": false,
	})
}

type checkDuplicateAssetReq struct {
	DeviceAssetId string `json:"deviceAssetId" binding:"required"`
	DeviceId      string `json:"deviceId" binding:"required"`
}

func CheckDuplicateAsset(c *gin.Context) {
	user := c.MustGet("user").(*view.User)

	req := checkDuplicateAssetReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	assetID, err := model.FindAssetByAssetIDAndDeviceID(user.ID, req.DeviceId, req.DeviceAssetId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var jsonBody gin.H
	if assetID != nil {
		jsonBody = gin.H{"isExist": true, "id": *assetID}
	} else {
		jsonBody = gin.H{"isExist": false}
	}
	c.JSON(http.StatusOK, jsonBody)
}
