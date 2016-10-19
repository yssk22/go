package xlog

import (
	"fmt"
	"os"
	"time"

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
		`{{formattimestamp .Timestamp}} [{{.Level}}] [{{.Goroutine}}] {{.Data}}`,
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
	// 2016-01-01T12:10:25Z [info] [3] This is a log
	// 2016-01-01T12:10:25Z [info] [3] T{1}
}

func ExampleLoggerWithLevelFilter() {
	formatter := NewTextFormatter(
		`{{formattimestamp .Timestamp}} [{{.Level}}] [{{.Goroutine}}] {{.Data}}`,
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
	// 2016-01-01T12:10:25Z [info] [3] This is a log
}
