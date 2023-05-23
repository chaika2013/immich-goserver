package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                string `gorm:"unique;index"`
	Password             string
	FirstName            string
	LastName             string
	ShouldChangePassword bool
	IsAdmin              bool
}

type Asset struct {
	gorm.Model

	// asset belongs to user
	UserID uint `gorm:"index"`
	User   User

	// data from upload request
	AssetType      string
	DeviceID       string `gorm:"index"`
	DeviceAssetID  string // this is original file name + '-' + file size
	FileCreatedAt  string
	FileModifiedAt string
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

type Exif struct {
	gorm.Model
	AssetID         uint
	Make            *string    // OLYMPUS OPTICAL CO.,LTD
	CameraModel     *string    `gorm:"column:model"` // C740UZ
	ExifImageWidth  *int       // 2048
	ExifImageHeight *int       // 1536
	Orientation     *int       // Horizontal (normal)
	CreateDate      *time.Time // 2003:01:08 18:12:09
	ModifyDate      *time.Time // 2003:01:08 18:12:09
	LensModel       *string    // SMC Pentax A 35-70mm
	FNumber         *float32   // 2.8
	FocalLength     *float32   // 6.3 mm
	ISO             *int       // 100
	ExposureTime    *string    // 1/30
}
