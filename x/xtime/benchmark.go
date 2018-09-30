package xtime

import "time"

// Benchmark measure the time.Duration for func f execution.
func Benchmark(f func()) time.Duration {
	t := time.Now()
	f()
	return time.Now().Sub(t)
}
