package types

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"strconv"
)

var ColorEmpty = &Color{empty: true}

type Color struct {
	empty bool
	col   color.RGBA
}

func (c *Color) RGBA() color.RGBA {
	return color.RGBA(c.col)
}

func (c *Color) String() string {
	if c == nil {
		return "null"
	}

	rgba := c.RGBA()

	col := fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)

	return col
}

func (c *Color) Bytes() []byte {
	if c == nil {
		return []byte{}
	}

	rgba := c.RGBA()

	return []byte{rgba.R, rgba.G, rgba.B}
}

func ParseColor(rawColor string) (*Color, error) {
	if rawColor == "" || rawColor == "null" {
		return ColorEmpty, nil
	}

	if len(rawColor) != 7 || rawColor[0] != '#' {
		return nil, fmt.Errorf("invalid color format")
	}

	r, err := strconv.ParseUint(rawColor[1:3], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid color format: could not parse red value: %v", err)
	}
	g, err := strconv.ParseUint(rawColor[3:5], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid color format: could not parse green value: %v", err)
	}
	b, err := strconv.ParseUint(rawColor[5:7], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid color format: could not parse blue value: %v", err)
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
	col := &Color{
		empty: false,
		col:   rgba,
	}
	return col
}

func ColorFromBytes(bytes []byte) *Color {
	if bytes == nil || len(bytes) != 3 {
		return nil
	}

	rgba := color.RGBA{
		R: bytes[0],
		G: bytes[1],
		B: bytes[2],
		A: 255,
	}

	col := ColorFromRGBA(rgba)
	return col
}

func (c *Color) HSL() (int, int, int) {
	rgb := c.RGBA()

	r := float64(rgb.R) / 255
	g := float64(rgb.G) / 255
	b := float64(rgb.B) / 255

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2
	var h, s float64

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			if g < b {
				h = (g-b)/d + 6
			} else {
				h = (g - b) / d
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		default:
			h = 0
		}
	}

	h *= 60
	s *= 100
	l *= 100

	return int(math.Round(h)), int(math.Round(s)), int(math.Round(l))
}

func (c *Color) MarshalJSON() ([]byte, error) {
	if c.empty {
		return json.Marshal(nil)
	}

	return json.Marshal(c.String())
}

func (c *Color) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == "null" {
		*c = *ColorEmpty
		return nil
	}

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

func (c *Color) IsEmpty() bool {
	return c == nil || c.empty
}

func (c *Color) distance(other *Color) uint {
	if c == nil || other == nil {
		return ^uint(0)
	}

	h1, s1, l1 := c.HSL()
	h2, s2, l2 := other.HSL()

	dh := h1 - h2
	dh *= dh * 2

	ds := s1 - s2
	ds *= ds

	dl := l1 - l2
	dl *= dl

	return uint(dh + ds + dl)
}

func (c *Color) Equals(other *Color) bool {
	if c == nil || other == nil {
		return false
	}

	return c.RGBA() == other.RGBA()
}
