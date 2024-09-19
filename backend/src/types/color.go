package types

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
)

var ColorDefault = ColorFromVals(25, 25, 50)
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
		return ColorDefault.String()
	}

	rgba := c.RGBA()

	col := fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)

	return col
}

func (c *Color) Bytes() []byte {
	if c == nil {
		return ColorDefault.Bytes()
	}

	rgba := c.RGBA()

	return []byte{rgba.R, rgba.G, rgba.B}
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

func (c *Color) MarshalJSON() ([]byte, error) {
	if c.empty {
		//return json.Marshal(nil)
		return json.Marshal(ColorDefault)
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
