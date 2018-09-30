package xlog

import (
	"os"
	"time"

	"github.com/yssk22/go/ansi"
	"github.com/yssk22/go/x/xtime"
)

func ExampleFormatter() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}`,
	)
	formatter.SetError(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}{{formatstack .}}`,
		ansi.Reset,
	)
	logger := New(NewIOSinkWithFormatter(os.Stdout, formatter))
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			logger.Errorf("This is a log")
		},
	)
	// Do not run tests since line number may differ in each environment.
	// Sample output is like follows:
	//
	//  2016-01-01T12:10:25Z [error] This is a log
	// 	    github.com/yssk22/go/x/xlog.ExampleLoggerWithConfigurationByLevel.func1 (at github.com/yssk22/go/x/xlog/xlog_test.go#75)
	// 	    github.com/yssk22/go/x/xtime.RunAt (at github.com/yssk22/go/x/xtime/xtime.go#33)
	//  	github.com/yssk22/go/x/xlog.ExampleLoggerWithConfigurationByLevel (at github.com/yssk22/go/x/xlog/xlog_test.go#82)
}
