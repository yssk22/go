// Package xlog provides extended utility functions for log
package xlog

import (
	"fmt"
	"time"

	"context"

	"github.com/speedland/go/x/xruntime"
	"github.com/speedland/go/x/xtime"
)

// Record is a data set for one log line/data
type Record struct {
	Timestamp time.Time
	Data      interface{}

	Level Level // Log Levell

	LoggerKey interface{}       // Logger Key
	Stack     []*xruntime.Frame // Log Source BenchmarkLoggerFewStackCapture

	ctx context.Context // used for application context
}

// Context returns context.Context associated with the record.
func (r *Record) Context() context.Context {
	return r.ctx
}

// ContextValue returns the context value
func (r *Record) ContextValue(key interface{}) interface{} {
	if r.ctx != nil {
		return r.ctx.Value(key)
	}
	return nil
}

// Option is option fields for logger.
type Option struct {
	MinStackCaptureOn Level // Minimum Level where logger captures stacks. Strongly recommend to set LevelFatal.
	StackCaptureDepth int   // # of stack frames to be captured on logging.
}

// Logger the logger
type Logger struct {
	*Option
	key    interface{}
	sink   Sink
	ctx    context.Context
	prefix string
}

// New returns a new Logger instance.
func New(s Sink) *Logger {
	option := new(Option)
	*option = *defaultOption // copy default option
	caller := xruntime.CaptureCaller()
	return &Logger{
		Option: option,
		sink:   s,
		key:    caller.PackageName,
		ctx:    nil,
	}
}

// WithContext returns a shallow copy of Logger with its context changed to ctx.
// The provided ctx must be non-nil.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		panic("nil context")
	}
	l2 := new(Logger)
	*l2 = *l
	l2.ctx = ctx
	return l2
}

// WithKey returns a shallow copy of Logger with its key changed to `key`.
func (l *Logger) WithKey(key interface{}) *Logger {
	l2 := new(Logger)
	*l2 = *l
	l2.key = key
	return l2
}

// WithPrefix returns a shallow copy of Logger with its prefix changed to `prefix`.
func (l *Logger) WithPrefix(prefix string) *Logger {
	l2 := new(Logger)
	*l2 = *l
	l2.prefix = prefix
	return l2
}

// Tracef to write text log with LevelTrace
func (l *Logger) Tracef(s string, v ...interface{}) {
	l.write(LevelTrace, l.format(s, v...))
}

// Trace to write data log with LevelTrace
func (l *Logger) Trace(v interface{}) {
	l.write(LevelTrace, v)
}

// Debugf to write text log with LevelDebug
func (l *Logger) Debugf(s string, v ...interface{}) {
	l.write(LevelDebug, l.format(s, v...))
}

// Debug to write data log with LevelDebug
func (l *Logger) Debug(v interface{}) {
	l.write(LevelDebug, v)
}

// Infof to write text log with LevelInfo
func (l *Logger) Infof(s string, v ...interface{}) {
	l.write(LevelInfo, l.format(s, v...))
}

// Info to write data log with LevelInfo
func (l *Logger) Info(v interface{}) {
	l.write(LevelInfo, v)
}

// Warnf to write text log with LevelWarn
func (l *Logger) Warnf(s string, v ...interface{}) {
	l.write(LevelWarn, l.format(s, v...))
}

// Warn to write data log with LevelWarn
func (l *Logger) Warn(v interface{}) {
	l.write(LevelWarn, v)
}

// Errorf to write text log with LevelError
func (l *Logger) Errorf(s string, v ...interface{}) {
	l.write(LevelError, l.format(s, v...))
}

// Error to write data log with LevelError
func (l *Logger) Error(v interface{}) {
	l.write(LevelError, v)
}

// Fatalf to write text log with LevelTrace
func (l *Logger) Fatalf(s string, v ...interface{}) {
	l.write(LevelFatal, l.format(s, v...))
}

// Fatal to write text log with LevelFatal
func (l *Logger) Fatal(v interface{}) {
	l.write(LevelFatal, v)
}

func (l *Logger) format(s string, v ...interface{}) string {
	if l.prefix == "" {
		return fmt.Sprintf(s, v...)
	}
	return fmt.Sprintf("%s%s", l.prefix, fmt.Sprintf(s, v...))
}

func (l *Logger) write(level Level, data interface{}) {
	if v, ok := data.(func(*Printer)); ok {
		data = printerFunc(v)
	}
	r := &Record{
		Level:     level,
		Data:      data,
		Timestamp: xtime.Now(),
		LoggerKey: l.key,
		ctx:       l.ctx,
	}
	if l.MinStackCaptureOn <= r.Level {
		r.Stack = xruntime.CaptureStackFrom(2, l.StackCaptureDepth)
	}
	l.sink.Write(r)
}
