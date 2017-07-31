package xlog

import (
	"fmt"
	"os"
	"time"

	"github.com/speedland/go/x/xtime"
	"context"
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

func ExampleLogger_WithContext() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] [{{.ContextValue "id"}}] {{.Data}}`,
	)
	logger := New(
		NewIOSinkWithFormatter(os.Stdout, formatter),
	)
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			logger.Errorf("This is a log")
			ctx := context.Background()
			ctx = context.WithValue(ctx, "id", "123")
			logger = logger.WithContext(ctx)
			logger.Errorf("This is a log")
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [error] [<no value>] This is a log
	// 2016-01-01T12:10:25Z [error] [123] This is a log
	//
}
