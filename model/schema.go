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

	// user
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
	OriginalDateTime *time.Time `gorm:"index"` // filled in once exif is parsed
	InLibrary        bool       // false if asset is in upload path
	AssetPath        string     // file name within the current path

	// exif
	ExifID uint
	Exif   Exif
}

type Exif struct {
	gorm.Model
}
