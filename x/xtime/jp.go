package xtime

import (
	"regexp"
	"strconv"
	"time"
)

// JST timezone
var JST = time.FixedZone("JST", 9*60*60)

var jpDateRegexp = regexp.MustCompile(`(?:(\d{4})年)?(\d{1,2})月(\d{1,2})日`)

// ParseJPDate parse the string expression of JP date YYYY年MM月DD日 and returns it as time.Time
// YYYY年 can be omitted and hintYear is used for that case. If `hintYear` is 0, Now().Year() is used as an year.
func ParseJPDate(s string, hintYear int) time.Time {
	if matched := jpDateRegexp.Copy().FindStringSubmatch(s); matched != nil {
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
