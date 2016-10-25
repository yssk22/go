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
		LevelInfo.Filter(
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
