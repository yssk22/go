package retry

import (
	"fmt"
	"time"

	"context"
)

// Do retries `task` function until it returns nil error.
// If task return an error, checker checks if the retry is needed afeter waitng by waiter.
func Do(ctx context.Context, task func(context.Context) error, backoff Backoff, checker Checker) error {
	var err error
	var attempt int
	for {
		err = task(ctx)
		attempt++
		if err == nil {
			return nil
		}
		if !checker.NeedRetry(ctx, attempt, err) {
			return err
		}
		ticker := time.NewTicker(backoff.Calc(ctx, attempt))
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry canceled: %v", ctx.Err())
		case <-ticker.C:
			break
		}
	}
	return err
}

// Backoff is an interface to implement backoff algorithm
type Backoff interface {
	Calc(context.Context, int) time.Duration
}

// ConstBackoff return a Backoff to return a constant time.Duration backoff given by `t`
func ConstBackoff(t time.Duration) Backoff {
	return &constBackoff{
		interval: t,
	}
}

type constBackoff struct {
	interval time.Duration
}

// Calc implements Backoff#Calc
func (b *constBackoff) Calc(ctx context.Context, attempt int) time.Duration {
	return b.interval
}

// Checker is an interface to check if a retry is needed or not.
type Checker interface {
	NeedRetry(context.Context, int, error) bool
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

func (u *until) NeedRetry(ctx context.Context, attempt int, err error) bool {
	return time.Now().Before(u.t)
}

// MaxRetries returns a Checker to check the number of retries is less than max.
func MaxRetries(max int) Checker {
	return &maxRetries{
		max: max,
	}
}

type maxRetries struct {
	max int
}

func (mr *maxRetries) NeedRetry(ctx context.Context, attempt int, err error) bool {
	return attempt < mr.max
}
