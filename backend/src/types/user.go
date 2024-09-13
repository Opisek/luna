package types

type User struct {
	Id        ID     `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Algorithm string `json:"-"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
}
