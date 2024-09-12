package types

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Color *Color    `json:"color"`
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

func (u *Url) String() string {
	return u.URL().String()
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

func (id ID) Bytes() []byte {
	return []byte(id.String())
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

func (id ID) MarshalJSON() ([]byte, error) {
	if id == EmptyId() {
		return json.Marshal(nil)
	}
	return json.Marshal(id.String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var rawId string
	if err := json.Unmarshal(data, &rawId); err != nil {
		return err
	}
	ID, err := IdFromString(rawId)
	if err != nil {
		return err
	}
	*id = ID
	return nil
}

type Color color.RGBA

func (c *Color) RGBA() color.RGBA {
	return color.RGBA(*c)
}

func (c *Color) String() string {
	rgba := c.RGBA()

	col := fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)

	return col
}

func ParseColor(rawColor string) (*Color, error) {
	if len(rawColor) != 7 || rawColor[0] != '#' {
		return nil, fmt.Errorf("invalid color format")
	}

	r, err := strconv.ParseInt(rawColor[1:3], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid color format: could not parse red value: %v", err)
	}
	b, err := strconv.ParseInt(rawColor[3:5], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid color format: could not parse blue value: %v", err)
	}
	g, err := strconv.ParseInt(rawColor[5:7], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid color format: could not parse green value: %v", err)
	}

	return ColorFromVals(uint8(r), uint8(g), uint8(b)), nil
}

func ColorFromVals(r, g, b uint8) *Color {
	rgba := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}

	col := ColorFromRGBA(rgba)
	return col
}

func ColorFromRGBA(rgba color.RGBA) *Color {
	col := Color(rgba)
	return &col
}

func (c *Color) MarshalJSON() ([]byte, error) {
	if c == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(c.String())
}

func (c *Color) UnmarshalJSON(data []byte) error {
	var rawColor string
	if err := json.Unmarshal(data, &rawColor); err != nil {
		return err
	}
	color, err := ParseColor(rawColor)
	if err != nil {
		return err
	}
	*c = Color(*color)
	return nil
}

const (
	SourceCaldav = "caldav"
	SourceIcal   = "ical"
)

const (
	AuthNone   = "none"
	AuthBasic  = "basic"
	AuthBearer = "bearer"
)

const (
	HashPlain = "plain"
)

type PgxScanner interface {
	Scan(dest ...interface{}) (err error)
}
