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

// ParseInt64Or parse string and return int64 value. The invalid string returns `or` value.
func ParseInt64Or(s string, or int64) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return or
	}
	return i
}

// MustParseInt64 parse string and return int64 value. it panics if an invalid string is passed.
func MustParseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
