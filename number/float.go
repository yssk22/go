package number

import "strconv"

// ParseFloatOr parse string and return float64 value. The invalid string returns `or` value.
func ParseFloatOr(s string, or float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return or
	}
	return f
}

// MustParseFloat parse string and return float64 value. it panics if an invalid string is passed.
func MustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
