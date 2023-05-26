package view

import "time"

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	// ProfileImagePath     string    `json:"profileImagePath,omitempty"`
	ShouldChangePassword bool `json:"shouldChangePassword"`
	IsAdmin              bool `json:"isAdmin"`
	// OAuthID              string    `json:"oauthId,omitempty"`
}
