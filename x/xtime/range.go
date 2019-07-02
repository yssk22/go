package xtime

import "time"

const (
	rangeDay = 24 * time.Hour
)

// RangeDay returns the start time and end time in that day.
func RangeDay(t time.Time) (time.Time, time.Time) {
	base := time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		t.Location(),
	)
	return base, base.Add(rangeDay)
}

// RangeMonth returns the start time and end time in that month.
func RangeMonth(t time.Time) (time.Time, time.Time) {
	base := time.Date(
		t.Year(), t.Month(), 1,
		0, 0, 0, 0,
		t.Location(),
	)
	if t.Month() == time.December {
		return base, time.Date(
			t.Year()+1, 1, 1,
			0, 0, 0, 0,
			t.Location(),
		)
	}
	return base, time.Date(
		t.Year(), t.Month()+1, 1,
		0, 0, 0, 0,
		t.Location(),
	)
}
