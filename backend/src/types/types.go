package types

import (
	"encoding/json"
	"image/color"
	"net/url"
	"time"
)

type Calendar struct {
	Source string      `json:"-"`
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Desc   string      `json:"desc"`
	Color  *color.RGBA `json:"color"`
}

type Event struct {
	Name  string      `json:"name"`
	Start time.Time   `json:"start"`
	End   time.Time   `json:"end"`
	Color *color.RGBA `json:"color"`
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Algorithm string `json:"-"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
}

type Url url.URL

func (u *Url) MarshalJSON() ([]byte, error) {
	if u == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(u.URL().String())
}

func (u *Url) UnmarshalJSON(data []byte) error {
	var rawUrl string
	if err := json.Unmarshal(data, &rawUrl); err != nil {
		return err
	}
	URL, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	*u = Url(*URL)
	return nil
}

func (u *Url) URL() *url.URL {
	return (*url.URL)(u)
}

func NewUrl(rawUrl string) (*Url, error) {
	URL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	return (*Url)(URL), nil
}
