package xtime

import (
	"fmt"
	"strings"
	"time"
)

// Formatter provides format functions for time
type Formatter struct {
	DateSeperator    string
	TimeSeperator    string
	ShowSeconds      bool
	humanOffsetHours int
}

// DefaultFormatter is a default Formatter
var DefaultFormatter = &Formatter{
	DateSeperator:    "/",
	TimeSeperator:    ":",
	ShowSeconds:      false,
	humanOffsetHours: 0,
}

// Humanize returns a new formatter to use humanized offset date/hours to format.
func (f *Formatter) Humanize(offsetHours int) *Formatter {
	return &Formatter{
		DateSeperator:    f.DateSeperator,
		TimeSeperator:    f.TimeSeperator,
		ShowSeconds:      f.ShowSeconds,
		humanOffsetHours: offsetHours,
	}
}

// FormatDateString formats the date part of a time.Time to string.
func (f *Formatter) FormatDateString(t time.Time) string {
	y, m, d, _ := getHumanDateTime(t, f.humanOffsetHours)
	return fmt.Sprintf("%d%s%02d%s%02d", y, f.DateSeperator, m, f.DateSeperator, d)
}

// FormatTimeString formats the time part of a time.Time to string.
func (f *Formatter) FormatTimeString(t time.Time) string {
	var s []string
	_, _, _, h := getHumanDateTime(t, f.humanOffsetHours)
	s = append(s, fmt.Sprintf("%02d", h))
	s = append(s, fmt.Sprintf("%02d", t.Minute()))
	if f.ShowSeconds {
		s = append(s, fmt.Sprintf("%02d", t.Second()))
	}
	return strings.Join(s, f.TimeSeperator)
}

// FormatDateTimeString formats the date and time part of a time.Time to string..
func (f *Formatter) FormatDateTimeString(t time.Time) string {
	return fmt.Sprintf("%s %s", f.FormatDateString(t), f.FormatTimeString(t))
}

// FormatDateString is a shortcut for DefaultFormatter.FormatDateString
func FormatDateString(t time.Time) string {
	return DefaultFormatter.FormatDateString(t)
}

// FormatTimeString is a shortcut for DefaultFormatter.FormatTimeString
func FormatTimeString(t time.Time) string {
	return DefaultFormatter.FormatTimeString(t)
}

// FormatDateTimeString is a shortcut for DefaultFormatter.FormatDateTimeString
func FormatDateTimeString(t time.Time) string {
	return DefaultFormatter.FormatDateTimeString(t)
}
