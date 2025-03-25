package types

type User struct {
	Id             ID     `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Admin          bool   `json:"admin"`
	Searchable     bool   `json:"searchable"`
	ProfilePicture *Url   `json:"profile_picture"`
}
