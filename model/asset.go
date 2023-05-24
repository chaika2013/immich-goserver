package model

import (
	"path/filepath"
	"time"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/helper"
)

type TimeBuckets struct {
	Count   uint             `json:"totalCount"`
	Buckets []TimeBucketInfo `json:"buckets"`
}

type TimeBucketInfo struct {
	Count      uint   `json:"count"`
	TimeBucket string `json:"timeBucket"`
}

func GetTimeBuckets(user *User) (*TimeBuckets, error) {
	timeBuckets := TimeBuckets{}
	if err := DB.Model(&Asset{}).Select("strftime(\"%Y-%m-01T00:00:00.000Z\", date_time_original) as time_bucket, count(id) as count").Where("user_id = ?", user.ID).Group("time_bucket").Order("time_bucket desc").Find(&timeBuckets.Buckets).Error; err != nil {
		return nil, err
	}

	var count int64
	if err := DB.Model(&Asset{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
		return nil, err
	}
	return &timeBuckets, nil
}

type AssetInfo struct {
	ID               string     `json:"id"`
	Type             string     `json:"type"` // Possible values: [IMAGE, VIDEO, AUDIO, OTHER]
	DeviceAssetID    string     `json:"deviceAssetId"`
	OwnerID          string     `json:"ownerId"`
	DeviceID         string     `json:"deviceId"`
	OriginalPath     string     `json:"originalPath"`
	OriginalFileName string     `json:"originalFileName"`
	ResizePath       string     `json:"resizePath"`
	FileCreatedAt    string     `json:"fileCreatedAt" gorm:"column:date_time_original"`
	FileModifiedAt   string     `json:"fileModifiedAt"`
	UpdatedAt        string     `json:"updatedAt"`
	IsFavorite       bool       `json:"isFavorite"`
	IsArchived       bool       `json:"isArchived"`
	MimeType         string     `json:"mimeType"`
	Duration         string     `json:"duration"`
	WebpPath         string     `json:"webpPath"`
	EncodedVideoPath string     `json:"encodedVideoPath,omitempty"`
	ExifInfo         *ExifInfo  `json:"exifInfo,omitempty"`
	SmartInfo        *SmartInfo `json:"smartInfo,omitempty"`
	LivePhotoVideoID string     `json:"livePhotoVideoId,omitempty"`
	Tags             []TagInfo  `json:"tags,omitempty"`
}

type ExifInfo struct {
	FileSizeInByte   uint64    `json:"fileSizeInByte,omitempty"`
	Make             string    `json:"make,omitempty"`
	Model            string    `json:"model,omitempty"`
	ExifImageWidth   uint      `json:"exifImageWidth,omitempty"`
	ExifImageHeight  uint      `json:"exifImageHeight,omitempty"`
	Orientation      string    `json:"orientation,omitempty"`
	DateTimeOriginal time.Time `json:"dateTimeOriginal,omitempty"`
	ModifyDate       time.Time `json:"modifyDate,omitempty"`
	TimeZone         string    `json:"timeZone,omitempty"`
	LensModel        string    `json:"lensModel,omitempty"`
	FNumber          float32   `json:"fNumber,omitempty"`
	FocalLength      float32   `json:"focalLength,omitempty"`
	ISO              uint      `json:"iso,omitempty"`
	ExposureTime     string    `json:"exposureTime,omitempty"`
	Latitude         float32   `json:"latitude,omitempty"`
	Longitude        float32   `json:"longitude,omitempty"`
	City             string    `json:"city,omitempty"`
	State            string    `json:"state,omitempty"`
	Country          string    `json:"country,omitempty"`
	Description      string    `json:"description,omitempty"`
}

type SmartInfo struct {
	Tags    []string `json:"tags,omitempty"`
	Objects []string `json:"objects,omitempty"`
}

type TagInfo struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // Possible values: [OBJECT, FACE, CUSTOM]
	Name        string `json:"name"`
	UserID      string `json:"userId"`
	RenameTagID string `json:"renameTagId,omitempty"`
}

func GetAssetsByTimeBuckets(user *User, timeBuckets []string) (assets []AssetInfo, err error) {
	err = DB.Model(&Asset{}).Where("user_id = ? and strftime(\"%Y-%m-01T00:00:00.000Z\", date_time_original) IN ?", user.ID, timeBuckets).Order("date_time_original desc").Find(&assets).Error
	return
}

func GetAssetIDsByDeviceID(user *User, deviceID string) (assetIDs []string, err error) {
	err = DB.Model(&Asset{}).Select("device_asset_id").Where("user_id = ? and device_id = ?", user.ID, deviceID).Find(&assetIDs).Error
	return
}

type UploadFile struct {
	AssetType      string `form:"assetType" binding:"required"`
	DeviceAssetID  string `form:"deviceAssetId" binding:"required"`
	DeviceID       string `form:"deviceId" binding:"required"`
	FileCreatedAt  string `form:"fileCreatedAt" binding:"required"`
	FileModifiedAt string `form:"fileModifiedAt" binding:"required"`
	IsFavorite     bool   `form:"isFavorite"`
	IsArchived     bool   `form:"isArchived"`
	IsVisible      bool   `form:"isVisible"`
	FileExtension  string `form:"fileExtension" binding:"required"`
	Duration       string `form:"duration"`
}

type UploadedAsset struct {
	ID        uint `json:"id"`
	Duplicate bool `json:"duplicate"`
}

func NewUploadAsset(user *User, uploadFile *UploadFile, originalFileName string,
	fileSize int64, crc32 uint32, fileName string) (*Asset, error) {

	asset := Asset{
		UserID: user.ID,

		AssetType:      uploadFile.AssetType,
		DeviceID:       uploadFile.DeviceID,
		DeviceAssetID:  uploadFile.DeviceAssetID,
		FileCreatedAt:  uploadFile.FileCreatedAt,
		FileModifiedAt: uploadFile.FileModifiedAt,
		IsFavorite:     uploadFile.IsFavorite,
		IsArchived:     uploadFile.IsArchived,
		IsVisible:      uploadFile.IsVisible,
		Duration:       uploadFile.Duration,

		OriginalFileName: originalFileName,
		FileSize:         fileSize,
		CRC32:            crc32,

		AssetPath: fileName,
	}

	// create asset and check for duplicates
	if err := DB.Create(&asset).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func GetAssetWithRealPathByID(assetID uint) (*Asset, error) {
	var asset Asset
	if err := DB.Find(&asset, assetID).Error; err != nil {
		return nil, err
	}

	basePath := config.UploadPath
	if asset.InLibrary {
		basePath = config.LibraryPath
	}
	asset.AssetPath = filepath.Join(*basePath, helper.StringID(asset.UserID), asset.AssetPath)
	return &asset, nil
}

func MoveAssetToLibrary(asset *Asset, newAssetPath string) error {
	return DB.Model(asset).Update("asset_path", newAssetPath).Error
}
