package xlog

import (
	"context"
	"os"
	"time"

	"github.com/speedland/go/x/xtime"
)

func Example_WithContext() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}`,
	)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			_defaultLogger := defaultLogger
			defer func() {
				defaultLogger = _defaultLogger
			}()
			defaultLogger = New(
				NewIOSinkWithFormatter(os.Stdout, formatter),
			)
			ctx, logger := WithContext(context.Background(), "[Context1] ")
			logger.Infof("This is a log")
			ctx, logger = WithContext(ctx, "[Context2] ")
			logger.Infof("This is a log")
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [info] [Context1] This is a log
	// 2016-01-01T12:10:25Z [info] [Context1] [Context2] This is a log
	//
}

func Example_WithContextAndKey() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.LoggerKey}}] [{{.Level}}] {{.Data}}`,
	)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			_defaultLogger := defaultLogger
			defer func() {
				defaultLogger = _defaultLogger
			}()
			defaultLogger = New(
				NewIOSinkWithFormatter(os.Stdout, formatter),
			)
			ctx, logger := WithContextAndKey(context.Background(), "[Context1] ", "MyKey1")
			logger.Infof("This is a log")
			ctx, logger = WithContextAndKey(ctx, "[Context2] ", "MyKey2")
			logger.Infof("This is a log")
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [MyKey1] [info] [Context1] This is a log
	// 2016-01-01T12:10:25Z [MyKey2] [info] [Context1] [Context2] This is a log
	//
}
