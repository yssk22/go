// Package benchmark provides a benchmark utility
package benchmark

import "time"

// Run measure the time to execute f.
func Run(f func() error) (time.Duration, error) {
	t1 := time.Now()
	err := f()
	t2 := time.Now()
	return t2.Sub(t1), err
}
