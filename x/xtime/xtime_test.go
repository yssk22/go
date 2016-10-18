package xtime

import (
	"fmt"
	"testing"
	"time"
	"github.com/speedland/go/x/testing/assert"
)

func ExampleRunAt() {
	t := time.Date(
		2016, 1, 1, 13, 12, 0, 0, time.UTC,
	)
	RunAt(t, func() {
		fmt.Println(Now())
		fmt.Println(Today())
	})
	// Output:
	// 2016-01-01 13:12:00 +0000 UTC
	// 2016-01-01 00:00:00 +0000 UTC
}

func Test_WaitAndEnsureAfter(t *testing.T) {
	assert := assert.New(t)
	base := Now()
	WaitAndEnsureAfter(base, 1*time.Second)
	assert.OK(Now().Sub(base) >= 1*time.Second)
}
