package xtime

import (
	"testing"
	"github.com/speedland/go/x/testing/assert"
)

func TestParseJPDate(t *testing.T) {
	a := assert.New(t)
	d := ParseJPDate("2016年6月12日", 0)
	a.EqInt(2016, d.Year())
	a.EqInt(6, int(d.Month()))
	a.EqInt(12, d.Day())
}

func TestParseJPDate_Short(t *testing.T) {
	a := assert.New(t)
	d := ParseJPDate("6月12日", 2010)
	a.EqInt(2010, d.Year())
	a.EqInt(6, int(d.Month()))
	a.EqInt(12, d.Day())
}
