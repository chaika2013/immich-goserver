package model

import (
	"time"

	"github.com/chaika2013/immich-goserver/view"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Email                string `gorm:"unique;index"`
	Password             string
	FirstName            string
	LastName             string
	ShouldChangePassword bool
	IsAdmin              bool
}

func GetUserCount(countAdminsOnly bool) (uint, error) {
	var count int64
	query := DB.Model(&User{})
	if countAdminsOnly {
		query = query.Where("is_admin = ?", true)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return uint(count), nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func GetUserByID(id uint) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	return &user, err
}

func GetUserByAPIKey(apiKey string) (*User, error) {
	// TODO
	var user User
	err := DB.First(&user).Error
	return &user, err
}

func GetAllUsers(userID uint, isAll bool) (users []view.User, err error) {
	db := DB.Model(&User{})
	if !isAll {
		db = db.Where("id <> ?", userID)
	}
	err = db.Find(&users).Error
	return
}
