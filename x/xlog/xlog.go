// Package xlog provides extended utility functions for log
package xlog

import (
	"fmt"
	"runtime"
	"time"

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
	Level     Level
	Data      interface{}

	SourceFile  string   // file path that generates the record
	SourceLine  int      // line number that generates the record
	SourceStack []string // Stack trace
	Goroutine   int      // Goroutine number
}

// Logger the logger
type Logger struct {
	sink Sink
}

// New returns a new Logger
func New(s Sink) *Logger {
	return &Logger{
		sink: s,
	}
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
	}
	// Set sourcefile and sourceline info.
	// This function would be called by public API like Trace(), Debug(),...
	// so that 2 levels upper caller would be actual file and line.
	_, file, line, ok := runtime.Caller(2)
	if ok {
		r.SourceFile = file
		r.SourceLine = line
	} else {
		// OK for testing not covered here
		r.SourceFile = "unknown"
		r.SourceLine = 0
	}
	r.Goroutine = runtime.NumGoroutine()
	// TODO: capture Stack
	l.sink.Write(r)
}
