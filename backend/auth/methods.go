package auth

import "net/http"

type AuthMethod interface {
	Do(req *http.Request) (*http.Response, error)
}

// No Authentication

type NoAuth struct{}

func (auth NoAuth) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func NewNoAuth() AuthMethod {
	return NoAuth{}
}

// Password and Username

type BasicAuth struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (auth BasicAuth) Do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(auth.Username, auth.Password)
	return http.DefaultClient.Do(req)
}

func NewBasicAuth(username, password string) AuthMethod {
	return BasicAuth{Username: username, Password: password}
}

// Bearer Token

type BearerAuth struct {
	Token string `json:"token" form:"token"`
}

func (auth BearerAuth) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+auth.Token)
	return http.DefaultClient.Do(req)
}

func NewBearerAuth(token string) AuthMethod {
	return BearerAuth{Token: token}
}
