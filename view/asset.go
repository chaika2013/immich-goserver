package view

import "time"

// time buckets
type TimeBuckets struct {
	Count   uint             `json:"totalCount"`
	Buckets []TimeBucketInfo `json:"buckets"`
}

type TimeBucketInfo struct {
	Count      uint   `json:"count"`
	TimeBucket string `json:"timeBucket"`
}

// asset list
type AssetInfo struct {
	ID               string `json:"id"`
	DeviceAssetID    string `json:"deviceAssetId"`
	OwnerID          string `json:"ownerId"`
	DeviceID         string `json:"deviceId"`
	Type             string `json:"type"` // Possible values: [IMAGE, VIDEO, AUDIO, OTHER]
	OriginalPath     string `json:"originalPath"`
	OriginalFileName string `json:"originalFileName"`
	ResizePath       string `json:"resizePath,omitempty"` // should be something if thumbnail is present
	FileCreatedAt    string `json:"fileCreatedAt,omitempty"`
	// FileModifiedAt string `json:"fileModifiedAt"`
	// UpdatedAt      string `json:"updatedAt"`
	IsFavorite bool `json:"isFavorite"`
	IsArchived bool `json:"isArchived"`
	// MimeType         string `json:"mimeType"`
	WebpPath string `json:"webpPath,omitempty"` // should be something if thumbnail is present
	// EncodedVideoPath string `json:"encodedVideoPath,omitempty"`
	Duration string `json:"duration"`
	// LivePhotoVideoID string `json:"livePhotoVideoId,omitempty"`
	// ExifInfo         *ExifInfo  `json:"exifInfo,omitempty"`
	// SmartInfo        *SmartInfo `json:"smartInfo,omitempty"`
	// Tags             []TagInfo  `json:"tags,omitempty"`
}

// type ExifInfo struct {
// 	FileSizeInByte   uint64    `json:"fileSizeInByte,omitempty"`
// 	Make             string    `json:"make,omitempty"`
// 	Model            string    `json:"model,omitempty"`
// 	ExifImageWidth   uint      `json:"exifImageWidth,omitempty"`
// 	ExifImageHeight  uint      `json:"exifImageHeight,omitempty"`
// 	Orientation      string    `json:"orientation,omitempty"`
// 	DateTimeOriginal time.Time `json:"dateTimeOriginal,omitempty"`
// 	ModifyDate       time.Time `json:"modifyDate,omitempty"`
// 	TimeZone         string    `json:"timeZone,omitempty"`
// 	LensModel        string    `json:"lensModel,omitempty"`
// 	FNumber          float32   `json:"fNumber,omitempty"`
// 	FocalLength      float32   `json:"focalLength,omitempty"`
// 	ISO              uint      `json:"iso,omitempty"`
// 	ExposureTime     string    `json:"exposureTime,omitempty"`
// 	Latitude         float32   `json:"latitude,omitempty"`
// 	Longitude        float32   `json:"longitude,omitempty"`
// 	City             string    `json:"city,omitempty"`
// 	State            string    `json:"state,omitempty"`
// 	Country          string    `json:"country,omitempty"`
// 	Description      string    `json:"description,omitempty"`
// }

// type SmartInfo struct {
// 	Tags    []string `json:"tags,omitempty"`
// 	Objects []string `json:"objects,omitempty"`
// }

// type TagInfo struct {
// 	ID          string `json:"id"`
// 	Type        string `json:"type"` // Possible values: [OBJECT, FACE, CUSTOM]
// 	Name        string `json:"name"`
// 	UserID      string `json:"userId"`
// 	RenameTagID string `json:"renameTagId,omitempty"`
// }

// uploading asset
type UploadFile struct {
	AssetType      string    `form:"assetType" binding:"required"`
	DeviceAssetID  string    `form:"deviceAssetId" binding:"required"`
	DeviceID       string    `form:"deviceId" binding:"required"`
	FileCreatedAt  time.Time `form:"fileCreatedAt" binding:"required"`
	FileModifiedAt time.Time `form:"fileModifiedAt" binding:"required"`
	IsFavorite     bool      `form:"isFavorite"`
	IsArchived     bool      `form:"isArchived"`
	IsVisible      bool      `form:"isVisible"`
	FileExtension  string    `form:"fileExtension" binding:"required"`
	Duration       string    `form:"duration"`
}

type UploadedAsset struct {
	ID        uint `json:"id"`
	Duplicate bool `json:"duplicate"`
}
