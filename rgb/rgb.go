package rgb

import (
	"fmt"
	"strconv"
)

// RGB is a alias for 32bit color interger
type RGB int

// ToHexString to get the hex string expression for the color
func (rgb RGB) ToHexString() string {
	r := (rgb & 0xFF0000) >> 16
	g := (rgb & 0x00FF00) >> 8
	b := (rgb & 0x0000FF)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// ParseRGB to parse the string and return as RGB
func ParseRGB(s string) (RGB, error) {
	l := len(s)
	if l == 6 {
		i, err := strconv.ParseInt(fmt.Sprintf("0x%s", s), 0, 32)
		return RGB(i), err
	} else if l == 7 && s[0] == '#' {
		i, err := strconv.ParseInt(fmt.Sprintf("0x%s", s[1:]), 0, 32)
		return RGB(i), err
	} else if l == 8 && s[0] == '0' && s[1] == 'x' {
		i, err := strconv.ParseInt(s, 0, 32)
		return RGB(i), err
	}
	return RGB(0), fmt.Errorf("Could not parse RGB color: %q", s)
}

// MustParseRGB is like ParseRGB but panic if an error occurrs
func MustParseRGB(s string) RGB {
	rgb, err := ParseRGB(s)
	if err != nil {
		panic(err)
	}
	return rgb
}

// MarshalJSON implements json.Marshaler#MarshalJSON()
func (rgb RGB) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", rgb.ToHexString())), nil
}

// UnmarshalJSON implements json.Unmarshaler#UnmarshalJSON()
func (rgb *RGB) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("Invalid string")
	}
	newval, err := ParseRGB(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*rgb = newval
	return nil
}
