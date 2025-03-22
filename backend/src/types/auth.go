package types

import "net/http"

type AuthMethod interface {
	Do(req *http.Request) (*http.Response, error)
	GetType() string
	String() (string, error)
}

type PasswordEntry struct {
	Hash       []byte
	Salt       []byte
	Algorithm  string
	Parameters map[string]int
}
