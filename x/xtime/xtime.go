// Package xtime provides extended utility functions for time
package xtime

import "time"

var nowFunc = time.Now

// Now returns time.Now() but you can set fixed return value by RunAt function
func Now() time.Time {
	return nowFunc()
}

// Today returns time.Today() but you can set fixed return value by RunAt function
func Today() time.Time {
	now := nowFunc()
	return time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
}

// RunAt sets the base time for Now().
// This is not concurrently safe so don't use this for parallel tests.
func RunAt(t time.Time, f func()) {
	nowFunc = func() time.Time {
		return t
	}
	defer (func() {
		nowFunc = time.Now
	})()
	f()
}

// WaitAndEnsureAfter waits for several time and ensure the duration passes after the base time.
func WaitAndEnsureAfter(base time.Time, d time.Duration) {
	diff := base.Add(d).Sub(Now())
	if diff > 0 {
		time.Sleep(diff)
		return
	}
}
