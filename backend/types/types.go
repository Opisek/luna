package types

import (
	"image/color"
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
