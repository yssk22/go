package ent

import (
	"strings"
	"time"

	"github.com/speedland/go/number"
	"github.com/speedland/go/x/xstrings"
	"github.com/speedland/go/x/xtime"
)

// ParseBool parse bool value
func ParseBool(v string) bool {
	norm := strings.ToLower(v)
	return norm == "true" || norm == "1"
}

// ParseInt parse int value
func ParseInt(v string) int {
	return number.MustParseInt(v)
}

// ParseInt64 parse int64 value
func ParseInt64(v string) int64 {
	return number.MustParseInt64(v)
}

// ParseFloat32 parse int64 value
func ParseFloat32(v string) float32 {
	return number.MustParseFloat32(v)
}

// ParseFloat64 parse int64 value
func ParseFloat64(v string) float64 {
	return number.MustParseFloat64(v)
}

// ParseStringList parse []string value
func ParseStringList(v string) []string {
	return xstrings.SplitAndTrim(v, ",")
}

// ParseTime parse time.Time value
func ParseTime(v string) time.Time {
	if v == "$now" {
		return xtime.Now()
	} else if v == "$today" {
		return xtime.Today()
	}
	return xtime.MustParse(v)
}
