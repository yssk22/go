package xtime

import (
	"regexp"
	"strconv"
	"time"
)

var dateRegxp = regexp.MustCompile(`(?:(\d{4})[\/-])?(\d{1,2})[\/-](\d{1,2})`)

// ParseDate parse the string expression of YYYY/MM/DD formatted time and returns it as time.Time
// YYYY/ can be omitted and hintYear is used for that case. The delimiter can be either '/' or '-'
func ParseDate(s string, location *time.Location, hintYear int) time.Time {
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
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, JST)
	}
	return time.Time{}
}

var timeRegexp = regexp.MustCompile(`(\d{1,2}):(\d{1,2})(?::(\d{1,2}))?`)

// ParseTime parse the string expression of HH:MM:SS formatted time and returns it as time.Time
// :SS can be omitted and hintYear is used for that case
func ParseTime(s string, location *time.Location) time.Time {
	if matched := timeRegexp.Copy().FindStringSubmatch(s); matched != nil {
		var h, m, s int
		h, _ = strconv.Atoi(matched[1])
		m, _ = strconv.Atoi(matched[2])
		if matched[3] != "" {
			s, _ = strconv.Atoi(matched[3])
		}
		return time.Date(0, 0, 0, h, m, s, 0, location)
	}
	return time.Time{}
}
