package xtime

import (
	"fmt"
	"time"
)

func ExampleFormatDateTimeString() {
	t := time.Date(
		2015, 5, 4, 2, 10, 0, 0, time.UTC,
	)
	fmt.Println(FormatDateTimeString(t))
	// Output:
	// 2015/05/04 02:10
}

func ExampleFormatter_Humanize() {
	t := time.Date(
		2015, 5, 4, 2, 10, 0, 0, time.UTC,
	)
	fmt.Println(DefaultFormatter.Humanize(4).FormatDateTimeString(t))
	// Output:
	// 2015/05/03 26:10
}
