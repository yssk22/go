package xtime

import (
	"testing"
	"time"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestRangeDay(t *testing.T) {
	a := assert.New(t)
	d, _ := Parse("2016-06-12T02:12:00Z")
	s, e := RangeDay(d)
	a.EqTime(
		time.Date(2016, 6, 12, 0, 0, 0, 0, time.UTC),
		s,
	)
	a.EqTime(
		time.Date(2016, 6, 13, 0, 0, 0, 0, time.UTC),
		e,
	)
}

func TestRangeMonth(t *testing.T) {
	a := assert.New(t)
	d, _ := Parse("2016-12-12T02:12:00Z")
	s, e := RangeMonth(d)
	a.EqTime(
		time.Date(2016, 12, 1, 0, 0, 0, 0, time.UTC),
		s,
	)
	a.EqTime(
		time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
		e,
	)
}
