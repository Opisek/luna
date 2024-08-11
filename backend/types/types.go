package types

import (
	"image/color"
)

type Calendar struct {
	Name  string      `json:"name"`
	Desc  string      `json:"desc"`
	Color *color.RGBA `json:"color"`
}

type Event struct {
	Name string
}
