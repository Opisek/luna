package types

import (
	"image/color"
)

type Calendar struct {
	Url   string      `json:"url"`
	Name  string      `json:"name"`
	Color *color.RGBA `json:"color"`
}

type Event struct {
	Name string
}
