package xtime

import (
	"testing"
	"time"
	"x/testing/assert"
)

func TestParseDate(t *testing.T) {
	a := assert.New(t)
	d := ParseDate("2016/6/12", time.UTC, 0)
	a.EqInt(2016, d.Year())
	a.EqInt(6, int(d.Month()))
	a.EqInt(12, d.Day())
}

func TestParseDate_Short(t *testing.T) {
	a := assert.New(t)
	d := ParseDate("6/12", time.UTC, 2010)
	a.EqInt(2010, d.Year())
	a.EqInt(6, int(d.Month()))
	a.EqInt(12, d.Day())
}

func TestParseTime(t *testing.T) {
	a := assert.New(t)
	d := ParseTime("6:12", JST)
	a.EqInt(6, d.Hour())
	a.EqInt(12, d.Minute())
}

func TestParseTime_Long(t *testing.T) {
	a := assert.New(t)
	d := ParseTime("6:12:13", JST)
	a.EqInt(6, d.Hour())
	a.EqInt(12, d.Minute())
	a.EqInt(13, d.Second())
}
