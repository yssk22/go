// Package xlog provides extended utility functions for log
package xlog

import (
	"golang.org/x/net/context"
)

var defaultLogger = New(NullSink)

// SetOption sets the option for the global logger
func SetOption(o *Option) {
	defaultLogger.Option = o
}

// SetOutput sets the sink for the global logger
func SetOutput(s Sink) {
	defaultLogger.sink = s
}

// WithContext returns a shallow copy of global Logger with its context changed to ctx.
// The provided ctx must be non-nil.
func WithContext(ctx context.Context) *Logger {
	return defaultLogger.WithContext(ctx)
}

// WithName returns a shallow copy of global Logger with its name changed to ctx.
func WithName(name string) *Logger {
	return defaultLogger.WithName(name)
}
