package xlog

import (
	"fmt"
	"os"
	"time"

	"github.com/speedland/go/x/ansi"
	"github.com/speedland/go/x/xtime"
)

type T struct {
	F string
}

func (t *T) String() string {
	return fmt.Sprintf("T{%s}", t.F)
}

func ExampleLogger() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}`,
	)
	logger := New(
		NewIOSinkWithFormatter(os.Stdout, formatter),
	)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			logger.Infof("This is a log")
			logger.Info(&T{"1"})
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [info] This is a log
	// 2016-01-01T12:10:25Z [info] T{1}
}

func ExampleLoggerWithLevelFilter() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}`,
	)
	logger := New(
		LevelFilter(
			LevelInfo,
			NewIOSinkWithFormatter(os.Stdout, formatter),
		),
	)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			logger.Infof("This is a log")
			logger.Debugf("DEBUGDEBUGDEBUG")
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [info] This is a log
}

func ExampleLoggerWithConfigurationByLevel() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}`,
	)
	formatter.SetError(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}{{formatstack .}}`,
		ansi.Reset,
	)
	logger := New(
		NewIOSinkWithFormatter(os.Stdout, formatter),
	)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			logger.Errorf("This is a log")
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [error] This is a log
	//	github.com/speedland/go/x/xlog.ExampleLoggerWithConfigurationByLevel.func1 (at github.com/speedland/go/x/xlog/xlog_test.go#75)
	//	github.com/speedland/go/x/xtime.RunAt (at github.com/speedland/go/x/xtime/xtime.go#33)
	//	github.com/speedland/go/x/xlog.ExampleLoggerWithConfigurationByLevel (at github.com/speedland/go/x/xlog/xlog_test.go#82)
}
