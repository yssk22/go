package xlog

import (
	"os"
	"time"

	"github.com/speedland/go/x/xtime"
)

func ExampleLogger_withLevelFilter() {
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

func ExampleLogger_withNameFilter() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] [{{.Name}}] {{.Data}}`,
	)
	// foo logger should output with the lever over LevelInfo.
	sink := NameFilter(
		map[string]Level{
			"foo": LevelInfo,
		},
		NewIOSinkWithFormatter(os.Stdout, formatter),
	)
	logger := New(sink)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			logger.Debugf("This should be written")
			logger = logger.WithName("foo")
			logger.Debugf("This should not be written")
			logger.Infof("This should be written")
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [debug] [github.com/speedland/go/x/xlog] This should be written
	// 2016-01-01T12:10:25Z [info] [foo] This should be written
}
