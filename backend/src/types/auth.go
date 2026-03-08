package types

import (
	"luna-backend/errors"
	"net/http"
)

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuthMethod interface {
	Do(req *http.Request) (*http.Response, *errors.ErrorTrace)
	GetType() string
	String() (string, error)
	HttpClient() HttpClientInterface
}

type PasswordEntry struct {
	Hash       []byte
	Salt       []byte
	Algorithm  string
	Parameters map[string]int
}
