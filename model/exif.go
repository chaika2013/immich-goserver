package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Exif struct {
	AssetID   uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Make            *string    // OLYMPUS OPTICAL CO.,LTD
	CameraModel     *string    `gorm:"column:model"` // C740UZ
	ExifImageWidth  *int       // 2048
	ExifImageHeight *int       // 1536
	Orientation     *int       // Horizontal (normal)
	CreateDate      *time.Time // 2003:01:08 18:12:09
	ModifyDate      *time.Time // 2003:01:08 18:12:09
	LensModel       *string    // SMC Pentax A 35-70mm
	FNumber         *float64   // 2.8
	FocalLength     *float64   // 6.3 mm
	ISO             *int       // 100
	ExposureTime    *string    // 1/30
}

func UpsertExif(exif *Exif, dateTimeOriginal *time.Time) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// create or update exif
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&exif).Error; err != nil {
			return err
		}

		// update the original date for the main asset
		if err := tx.Model(&Asset{ID: exif.AssetID}).Update("date_time_original", dateTimeOriginal).Error; err != nil {
			return err
		}

		return nil
	})
}
