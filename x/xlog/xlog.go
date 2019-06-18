// Package xlog provides the logging framework for applications.
//
// Logger is a top level struct that writes any type of log data to Sink.
//
// * Name
//
// Each of Logger should have it's name to identify which log records are sent by which logger.
// By default, that name is automatically resolved by the package name.
//
// * Level and lazy formatting
//
// As other logging frameworks do, xlog supports logging level support.And we support lazy evaluated logging.
//
//    logger.Debug(func(fmt *Printer){
//       result := somethingHeavy()
//       fmt.Println(result)
//    })
//
// In this semantic, somethingHeavy() is a heavy workload and only executed
// if minimum filter level is under Debug. If the level is upper than Debug,
// the heavy workload is not executed so that the logging cost will be reduced.
//
// * Pipeline
//
// Sink is an interface that implements `Write(*Record) error` and it's easy to
// make a pipeline from a Sink to another Sink. Filter is a such kind of Sink that creats
// pipeline for filtering records.
//
// * Context
//
// Logger can be aware of context.Context (even we support "context")
// You can associate any context by *Logger.WithContext() and write it by `{{context key .}}`
//
package xlog

import (
	"fmt"
	"os"

	"context"

	"github.com/yssk22/go/x/xcontext"
)

var defaultOption = &Option{
	MinStackCaptureOn: LevelError,
	StackCaptureDepth: 30,
}

var defaultIOFormatter = NewTextFormatter(
	`{{formattimestamp .}} [{{.Level}}] {{.Data}}{{formatstack .}}`,
)

var defaultFilter = LevelFilter(LevelInfo).Pipe(
	NewIOSinkWithFormatter(
		os.Stderr, defaultIOFormatter,
	),
)

var defaultLogger = New(defaultFilter)

// SetSink sets the sink for the global logger
func SetSink(s Sink) {
	defaultLogger = New(s)
}

// SetOption sets the option for the global logger
func SetOption(o *Option) {
	defaultOption = o
}

// WithKey returns a shallow copy of global Logger with its name changed to name.
func WithKey(name string) *Logger {
	return defaultLogger.WithKey(name)
}

var loggerContextKey = xcontext.NewKey("logger")

// WithContext returns a shallow copy of global Logger with its context changed to ctx.
func WithContext(ctx context.Context, prefix string) (context.Context, *Logger) {
	var instance *Logger
	if ctxLogger, ok := ctx.Value(loggerContextKey).(*Logger); ok {
		instance = new(Logger)
		*instance = *ctxLogger
		if prefix != "" {
			instance.prefix = fmt.Sprintf("%s%s", instance.prefix, prefix)
		}
	} else {
		instance = defaultLogger.WithContext(ctx)
		instance.prefix = prefix
	}
	newCtx := context.WithValue(ctx, loggerContextKey, instance)
	instance.ctx = newCtx
	return newCtx, instance
}

// WithContextAndKey returns a shallow copy of global Logger with its context changed to ctx and bound with `key`
func WithContextAndKey(ctx context.Context, prefix string, key interface{}) (context.Context, *Logger) {
	ctx, logger := WithContext(ctx, prefix)
	return ctx, logger.WithKey(key)
}
