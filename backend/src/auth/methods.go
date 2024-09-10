package auth

import (
	"encoding/json"
	"fmt"
	"luna-backend/types"
	"net/http"
)

type AuthMethod interface {
	Do(req *http.Request) (*http.Response, error)
	GetType() string
	String() (string, error)
}

// No Authentication

type NoAuth struct{}

func (auth NoAuth) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func (auth NoAuth) GetType() string {
	return types.AuthNone
}
func (auth NoAuth) String() (string, error) {
	return "", nil
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

func (auth BasicAuth) GetType() string {
	return types.AuthBasic
}
func (auth BasicAuth) String() (string, error) {
	bytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("could not marshal basic auth: %v", err)
	}
	return string(bytes), nil
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

func (auth BearerAuth) GetType() string {
	return types.AuthBearer
}
func (auth BearerAuth) String() (string, error) {
	bytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("could not marshal bearer auth: %v", err)
	}
	return string(bytes), nil
}

func NewBearerAuth(token string) AuthMethod {
	return BearerAuth{Token: token}
}
