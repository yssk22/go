package xtime

import (
	"testing"
	"time"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestSleepAndEnsureAfter(t *testing.T) {
	a := assert.New(t)
	base := Now()
	SleepAndEnsureAfter(base, 1*time.Second)
	a.OK(Now().Sub(base) >= 1*time.Second)
}
