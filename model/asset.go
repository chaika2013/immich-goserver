package model

import (
	"path/filepath"
	"time"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/helper"
	"github.com/chaika2013/immich-goserver/view"
	"gorm.io/gorm"
)

type Asset struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// asset belongs to user
	UserID uint `gorm:"index"`
	User   User

	// data from upload request
	AssetType      string
	DeviceID       string `gorm:"index"`
	DeviceAssetID  string // this is original file name + '-' + file size
	FileCreatedAt  time.Time
	FileModifiedAt time.Time
	IsFavorite     bool
	IsArchived     bool
	IsVisible      bool
	Duration       string `form:"duration"`

	// data calculated from the upload request
	OriginalFileName string
	FileSize         int64
	CRC32            uint32

	// info after asset was processed
	DateTimeOriginal *time.Time `gorm:"index"` // filled in once exif is parsed
	InLibrary        bool       // false if asset is in upload path
	AssetPath        string     // file name within the current path

	// has-one exif
	Exif Exif
}

func GetTimeBuckets(userID uint) (*view.TimeBuckets, error) {
	timeBuckets := view.TimeBuckets{}
	if err := DB.Model(&Asset{}).Select("strftime(\"%Y-%m-01T00:00:00.000Z\", date_time_original) as time_bucket, count(id) as count").Where("user_id = ?", userID).Group("time_bucket").Order("time_bucket desc").Find(&timeBuckets.Buckets).Error; err != nil {
		return nil, err
	}

	var count int64
	if err := DB.Model(&Asset{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, err
	}
	timeBuckets.Count = uint(count)

	return &timeBuckets, nil
}

func GetAssetsByTimeBuckets(userID uint, timeBuckets []string) (assets []view.AssetInfo, err error) {
	withEmptyBucket := false
	for _, bucket := range timeBuckets {
		if bucket == "" {
			withEmptyBucket = true
			break
		}
	}
	checkIsNull := ""
	if withEmptyBucket {
		checkIsNull = " or date_time_original is null"
	}
	query := "user_id = ? and (strftime(\"%Y-%m-01T00:00:00.000Z\", date_time_original) IN ?" + checkIsNull + ")"
	err = DB.Model(&Asset{}).Where(query, userID, timeBuckets).Order("date_time_original desc").Find(&assets).Error
	return
}

func GetAssetIDsByDeviceID(userID uint, deviceID string) (assetIDs []string, err error) {
	err = DB.Model(&Asset{}).Select("device_asset_id").Where("user_id = ? and device_id = ?", userID, deviceID).Find(&assetIDs).Error
	return
}

func NewUploadAsset(userID uint, uploadFile *view.UploadFile, originalFileName string,
	fileSize int64, crc32 uint32, fileName string) (*Asset, error) {

	asset := Asset{
		UserID: userID,

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
	return DB.Model(asset).Updates(Asset{AssetPath: newAssetPath, InLibrary: true}).Error
}

func FindAssetByAssetIDAndDeviceID(userID uint, deviceID string, deviceAssetID string) (*string, error) {
	var assetID string
	err := DB.Model(&Asset{}).Select("id").Where("device_id = ? and device_asset_id = ?", deviceID, deviceAssetID).First(&assetID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &assetID, nil
}
