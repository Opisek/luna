package types

import "time"

type User struct {
	Id             ID        `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Admin          bool      `json:"admin"`
	Verified       bool      `json:"verified"`
	Enabled        bool      `json:"enabled"`
	Searchable     bool      `json:"searchable"`
	ProfilePicture *Url      `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}
