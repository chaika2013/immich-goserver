package model

func GetUserCount(countAdminsOnly bool) uint {
	var count int64
	query := DB.Model(&User{})
	if countAdminsOnly {
		query = query.Where("is_admin = ?", true)
	}
	query.Count(&count)
	return uint(count)
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	r := DB.Where("email = ?", email).First(&user)
	return &user, r.Error
}

func GetUserByID(id uint) (*User, error) {
	var user User
	r := DB.First(&user, id)
	return &user, r.Error
}

func GetUserByAPIKey(apiKey string) (*User, error) {
	// TODO
	var user User
	r := DB.First(&user)
	return &user, r.Error
}
