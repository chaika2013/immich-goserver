package model

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
