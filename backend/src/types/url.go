package types

import (
	"encoding/json"
	"net/url"
)

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

func (u *Url) String() string {
	return u.URL().String()
}

func (u *Url) Subpage(subpages ...string) *Url {
	return (*Url)(u.URL().JoinPath(subpages...))
}

func (u *Url) Query() *url.Values {
	vals := u.URL().Query()
	return &vals
}

func (u *Url) SetQuery(query *url.Values) *Url {
	u.URL().RawQuery = query.Encode()
	return u
}

func NewUrl(rawUrl string) (*Url, error) {
	URL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	return (*Url)(URL), nil
}

func NewUrlSafe(rawUrl string) *Url {
	url, err := NewUrl(rawUrl)
	if err != nil {
		panic(err)
	}
	return url
}
