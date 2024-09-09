package types

import (
	"encoding/json"
	"image/color"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Calendar struct {
	Source ID          `json:"-"`
	Id     ID          `json:"id"`
	Path   string      `json:"path"`
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
	Id        ID     `json:"id"`
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

type ID uuid.UUID

func EmptyId() ID {
	return ID(uuid.Nil)
}

func RandomId() ID {
	id, _ := uuid.NewRandom()
	return IdFromUuid(id)
}

func (id ID) String() string {
	uuids := uuid.UUIDs([]uuid.UUID{uuid.UUID(id)})
	strings := uuids.Strings()
	return strings[0]
}

func (id ID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func IdFromBytes(bytes []byte) (ID, error) {
	id, err := uuid.FromBytes(bytes)
	if err != nil {
		return EmptyId(), err
	}
	return IdFromUuid(id), nil
}

func IdFromString(str string) (ID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return EmptyId(), err
	}
	return IdFromUuid(id), nil
}

func IdFromUuid(uu uuid.UUID) ID {
	return ID(uu)
}
