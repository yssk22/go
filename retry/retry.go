package retry

import (
	"time"

	"golang.org/x/net/context"
)

// Do retries `task` function until it returns nil error.
// If task return an error, checker checks if the retry is needed afeter waitng by waiter.
func Do(ctx context.Context, task func(context.Context) error, waiter Waiter, checker Checker) error {
	var err error
	for checker.Check() {
		if err = task(ctx); err == nil {
			return nil
		}
		waiter.Wait()
	}
	return err
}

// Waiter is an interface to wait for the next retry
type Waiter interface {
	Wait()
}

// Interval return a Waiter to wait for a given `t` interval.
func Interval(t time.Duration) Waiter {
	return &intervalWaiter{
		interval: t,
	}
}

// IntervalWaiter is a waiter to wait for the given interval time.
type intervalWaiter struct {
	interval time.Duration
}

// Wait implements Waiter#Wait()
func (w *intervalWaiter) Wait() {
	time.Sleep(w.interval)
}

// Checker is an interface to check if a retry is needed or not.
type Checker interface {
	Check() bool
}

// Until returns a Checker to check the time is before `t`
func Until(t time.Time) Checker {
	return &until{
		t: t,
	}
}

type until struct {
	t time.Time
}

func (u *until) Check() bool {
	return time.Now().Before(u.t)
}

// MaxRetries returns a Checker to check the number of retries is less than max.
func MaxRetries(max int) Checker {
	return &maxRetries{
		retries: 0,
		max:     max,
	}
}

type maxRetries struct {
	retries int
	max     int
}

func (mr *maxRetries) Check() bool {
	mr.retries++
	return mr.retries <= mr.max
}
