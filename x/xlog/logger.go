// Package xlog provides extended utility functions for log
package xlog

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/speedland/go/x/xruntime"
	"github.com/speedland/go/x/xtime"
)

// Level is a enum for log Level
//go:generate enum -type=Level
type Level int

// Available Level values
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Record is a data set for one log line/data
type Record struct {
	Timestamp time.Time
	Data      interface{}

	Level Level             // Log Levell
	Name  string            // Log Name
	Stack []*xruntime.Frame // Log Source Stack

	ctx context.Context // used for application context
}

// Option is option fields for logger.
type Option struct {
	MinStackCaptureOn Level // Minimum Level where logger captures stacks. Strongly recommend to set LevelFatal.
	StackCaptureDepth int   // # of stack frames to be captured on logging.
}

// Logger the logger
type Logger struct {
	*Option
	name string
	sink Sink
	ctx  context.Context
}

// New is like NewWithSender but name string is automatically set as the current package name.
func New(s Sink) *Logger {
	return NewWithName(xruntime.CaptureStackFrom(1, 1)[0].PackageName, s)
}

// NewWithName returns a new Logger with name.
func NewWithName(name string, s Sink) *Logger {
	const maxDepth = 50
	return &Logger{
		Option: &Option{
			MinStackCaptureOn: LevelError,
			StackCaptureDepth: maxDepth,
		},
		name: name,
		sink: s,
		ctx:  nil,
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

// WithName returns a shallow copy of Logger with its name changed to name.
func (l *Logger) WithName(name string) *Logger {
	l2 := new(Logger)
	*l2 = *l
	l2.name = name
	return l2
}

// Tracef to write text log with LevelTrace
func (l *Logger) Tracef(s string, v ...interface{}) {
	l.write(LevelTrace, fmt.Sprintf(s, v...))
}

// Trace to write data log with LevelTrace
func (l *Logger) Trace(v interface{}) {
	l.write(LevelTrace, v)
}

// Debugf to write text log with LevelDebug
func (l *Logger) Debugf(s string, v ...interface{}) {
	l.write(LevelDebug, fmt.Sprintf(s, v...))
}

// Debug to write data log with LevelDebug
func (l *Logger) Debug(v interface{}) {
	l.write(LevelDebug, v)
}

// Infof to write text log with LevelInfo
func (l *Logger) Infof(s string, v ...interface{}) {
	l.write(LevelInfo, fmt.Sprintf(s, v...))
}

// Info to write data log with LevelInfo
func (l *Logger) Info(v interface{}) {
	l.write(LevelInfo, v)
}

// Warnf to write text log with LevelWarn
func (l *Logger) Warnf(s string, v ...interface{}) {
	l.write(LevelError, fmt.Sprintf(s, v...))
}

// Warn to write data log with LevelWarn
func (l *Logger) Warn(v interface{}) {
	l.write(LevelWarn, v)
}

// Errorf to write text log with LevelError
func (l *Logger) Errorf(s string, v ...interface{}) {
	l.write(LevelError, fmt.Sprintf(s, v...))
}

// Error to write data log with LevelError
func (l *Logger) Error(v interface{}) {
	l.write(LevelWarn, v)
}

// Fatalf to write text log with LevelTrace
func (l *Logger) Fatalf(s string, v ...interface{}) {
	l.write(LevelFatal, fmt.Sprintf(s, v...))
}

// Fatal to write text log with LevelFatal
func (l *Logger) Fatal(v interface{}) {
	l.write(LevelFatal, v)
}

func (l *Logger) write(level Level, data interface{}) {
	r := &Record{
		Level:     level,
		Data:      data,
		Timestamp: xtime.Now(),
		Name:      l.name,
		ctx:       l.ctx,
	}
	if l.MinStackCaptureOn <= r.Level {
		r.Stack = xruntime.CaptureStackFrom(2, l.StackCaptureDepth)
	}
	l.sink.Write(r)
}
