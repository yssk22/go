package xtime

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleTimestampMarshalJSON() {
	t := Timestamp(time.Date(
		2015, 5, 4, 2, 10, 0, 0, time.UTC,
	))
	s, _ := json.Marshal(&t)
	fmt.Println(string(s))
	// Output:
	// 1430705400
}

func ExampleTimestampUnmarshalJSON() {
	var t Timestamp
	json.Unmarshal([]byte("1430705400"), &t)
	fmt.Println(time.Time(t).Unix())
	// Output:
	// 1430705400
}
