package xtime

import (
	"testing"
	"time"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestParse(t *testing.T) {
	a := assert.New(t)
	d, _ := Parse("2016-06-12T02:12:00Z")
	a.EqTime(
		time.Date(2016, 6, 12, 2, 12, 0, 0, time.UTC),
		d,
	)
}

func TestParseDate(t *testing.T) {
	a := assert.New(t)
	d, _ := ParseDate("2016/6/12", time.UTC, 0)
	a.EqInt(2016, d.Year())
	a.EqInt(6, int(d.Month()))
	a.EqInt(12, d.Day())
}

func TestParseDate_Short(t *testing.T) {
	a := assert.New(t)
	d, _ := ParseDate("6/12", time.UTC, 2010)
	a.EqInt(2010, d.Year())
	a.EqInt(6, int(d.Month()))
	a.EqInt(12, d.Day())
}

func TestParseTime(t *testing.T) {
	a := assert.New(t)
	d, _ := ParseTime("6:12", JST)
	a.EqInt(6, d.Hour())
	a.EqInt(12, d.Minute())
}

func TestParseTime_Long(t *testing.T) {
	a := assert.New(t)
	d, _ := ParseTime("6:12:13", JST)
	a.EqInt(6, d.Hour())
	a.EqInt(12, d.Minute())
	a.EqInt(13, d.Second())
}
