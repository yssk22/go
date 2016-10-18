package xtime

import "time"

// HumanToday returns Today for humans (considering time offset)
// The time between 0am - 4am should be considerd as the previous day if the offset is 4.
func HumanToday(offsetHours int) time.Time {
	now := Now()
	y, m, d, _ := getHumanDateTime(now, offsetHours)
	return time.Date(
		y, m, d, offsetHours, 0, 0, 0,
		now.Location(),
	)
}

// HumanTodayIn returns Today for humans (considering time offset) in the given location.
func HumanTodayIn(offsetHours int, loc *time.Location) time.Time {
	localNow := Now().In(loc)
	y, m, d, _ := getHumanDateTime(localNow, offsetHours)
	return time.Date(
		y, m, d, offsetHours, 0, 0, 0,
		localNow.Location(),
	)
}

func getHumanDateTime(t time.Time, offsetHours int) (y int, m time.Month, d int, h int) {
	y, m, d = t.Date()
	h = t.Hour()
	if t.Hour() < offsetHours {
		t = t.Add(-24 * time.Hour)
		y, m, d = t.Date()
		h = h + 24
		return y, m, d, h
	}
	return y, m, d, h
}
