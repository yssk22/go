package number

import "strconv"

// ParseIntOr parse string and return int value. The invalid string returns `or` value.
func ParseIntOr(s string, or int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return or
	}
	return i
}

// MustParseInt parse string and return int value. it panics if an invalid string is passed.
func MustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
