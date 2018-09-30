package xtime

import "time"

// SleepAndEnsureAfter sleeps for several time to ensure the duration passes after the base time
// for the routines that executes at a regular cycle.
func SleepAndEnsureAfter(base time.Time, d time.Duration) {
	diff := base.Add(d).Sub(Now())
	if diff > 0 {
		time.Sleep(diff)
		return
	}
}
