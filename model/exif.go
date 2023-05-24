package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
