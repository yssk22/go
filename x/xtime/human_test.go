package xtime

import (
	"fmt"
	"time"
)

func ExampleHumanToday() {
	t := time.Date(
		2016, 5, 4, 2, 10, 0, 0, time.UTC,
	)
	RunAt(t, func() {
		fmt.Println(HumanToday(4))
	})
	// Output:
	// 2016-05-03 04:00:00 +0000 UTC
}

func ExampleHumanTodayIn() {
	// Case if we want to get today of different location from system one.
	// Now returns 2016/05/04 02:10 am in UTC, which is 2016/05/04 11:10 am in JST
	// So human today for JST should return 2016/05/04 04:00, not 2016/05/03 04:00
	t := time.Date(
		2016, 5, 4, 2, 10, 0, 0, time.UTC,
	)
	RunAt(t, func() {
		fmt.Println(HumanTodayIn(4, JST))
	})
	// Output:
	// 2016-05-04 04:00:00 +0900 JST
}
