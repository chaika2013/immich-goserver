package model

import (
	"github.com/chaika2013/immich-goserver/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() (err error) {
	DB, err = gorm.Open(sqlite.Open(*config.DatabasePath), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return
	}

	if err = DB.AutoMigrate(&Asset{}, &User{}, &Exif{}); err != nil {
		return
	}

	// TODO for testing, remove later
	{
		var count int64
		DB.Model(&User{}).Count(&count)
		if count == 0 {
			user := User{
				Email:                "test.user@gmail.com",
				Password:             "$2a$14$rRKBPSc.syVWf3AqoIvdXOEvb5Dq83WlxaO.La1/30Gt5ysB.TFzS",
				FirstName:            "Test",
				LastName:             "User",
				ShouldChangePassword: false,
				IsAdmin:              true,
			}
			DB.Create(&user)
		}
	}
	// {
	// 	var count int64
	// 	DB.Model(&Asset{}).Count(&count)
	// 	if count == 0 {
	// 		ts := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	// 		for i := 0; i < 2000; i++ {
	// 			origin := ts
	// 			for j := 0; j < 1+rand.Intn(20); j++ {
	// 				asset := Asset{
	// 					UserID:           1,
	// 					DeviceID:         "CLI",
	// 					OriginalFileName: "filename",
	// 					// OriginalDateTime: &origin,
	// 				}
	// 				DB.Create(&asset)
	// 				origin = origin.Add(time.Hour)
	// 			}
	// 			ts = ts.Add(7 * 24 * time.Hour)
	// 		}
	// 	}
	// }

	return nil
}
