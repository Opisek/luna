package sources

import "net/http"

type SourceAuth interface {
	Do(req *http.Request) (*http.Response, error)
}

// No Authentication

type noAuth struct{}

func (auth noAuth) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func NoAuth() SourceAuth {
	return noAuth{}
}

// Password and Username

type basicAuth struct {
	Username string
	Password string
}

func (auth basicAuth) Do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(auth.Username, auth.Password)
	return http.DefaultClient.Do(req)
}

func BasicAuth(username, password string) SourceAuth {
	return basicAuth{Username: username, Password: password}
}

// Bearer Token

type bearerAuth struct {
	Token string
}

func (auth bearerAuth) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+auth.Token)
	return http.DefaultClient.Do(req)
}

func BearerAuth(token string) SourceAuth {
	return bearerAuth{Token: token}
}
