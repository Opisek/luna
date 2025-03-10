package util

import (
	"errors"
	"regexp"
	"strings"
)

var characterRegex = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
var urlRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,63}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])")

func IsValidUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 25 {
		return errors.New("username must be at most 25 characters long")
	}
	if !characterRegex.MatchString(username) {
		return errors.New("username can only contain letters and numbers")
	}
	return nil
}

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 50 {
		return errors.New("password must be at most 50 characters long")
	}
	return nil
}

func IsValidUrl(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("url must start with \"http://\" or \"https://\"")
	}
	if !urlRegex.MatchString(url) {
		return errors.New("the url contains illegal characters or is invalid")
	}
	return nil
}

func IsValidEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}
