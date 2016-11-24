package xtime

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var dateRegxp = regexp.MustCompile(`(?:(\d{4})[\/-])?(\d{1,2})[\/-](\d{1,2})`)

// Parse parse RFC3339 format by default.
func Parse(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

// MustParse is like Parse but panic if an error occurrs.
func MustParse(value string) time.Time {
	t, e := Parse(value)
	if e != nil {
		panic(e)
	}
	return t
}

// ParseDate parse the string expression of YYYY/MM/DD formatted time and returns it as time.Time
// YYYY/ can be omitted and hintYear is used for that case. The delimiter can be either '/' or '-'
func ParseDate(s string, location *time.Location, hintYear int) (time.Time, error) {
	if matched := dateRegxp.Copy().FindStringSubmatch(s); matched != nil {
		var y, m, d int
		if matched[1] == "" {
			y = hintYear
			if y == 0 {
				y = Now().Year()
			}
		} else {
			y, _ = strconv.Atoi(matched[1])
		}
		m, _ = strconv.Atoi(matched[2])
		d, _ = strconv.Atoi(matched[3])
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, location), nil
	}
	return time.Time{}, fmt.Errorf("invalid date format")
}

// MustParseDate is like ParseDate but panic if an error occurrs.
func MustParseDate(s string, location *time.Location, hintYear int) time.Time {
	t, e := ParseDate(s, location, hintYear)
	if e != nil {
		panic(e)
	}
	return t
}

var timeRegexp = regexp.MustCompile(`(\d{1,2}):(\d{1,2})(?::(\d{1,2}))?`)

// ParseTime parse the string expression of HH:MM:SS formatted time and returns it as time.Time
// :SS can be omitted and hintYear is used for that case
func ParseTime(s string, location *time.Location) (time.Time, error) {
	if matched := timeRegexp.Copy().FindStringSubmatch(s); matched != nil {
		var h, m, s int
		h, _ = strconv.Atoi(matched[1])
		m, _ = strconv.Atoi(matched[2])
		if matched[3] != "" {
			s, _ = strconv.Atoi(matched[3])
		}
		return time.Date(0, 0, 0, h, m, s, 0, location), nil
	}
	return time.Time{}, fmt.Errorf("invalid time format")
}

// MustParseTime is like ParseTime but panic if an error occurrs.
func MustParseTime(s string, location *time.Location) time.Time {
	t, e := ParseTime(s, location)
	if e != nil {
		panic(e)
	}
	return t
}
