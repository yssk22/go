package number

import "strconv"

// ParseFloat32Or parse string and return float64 value. The invalid string returns `or` value.
func ParseFloat32Or(s string, or float32) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return or
	}
	return float32(f)
}

// MustParseFloat32 parse string and return float64 value. it panics if an invalid string is passed.
func MustParseFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(f)
}

// ParseFloat64Or parse string and return float64 value. The invalid string returns `or` value.
func ParseFloat64Or(s string, or float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return or
	}
	return f
}

// MustParseFloat64 parse string and return float64 value. it panics if an invalid string is passed.
func MustParseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
