package xlog

import (
	"os"
	"testing"
	"time"

	"github.com/speedland/go/x/xtesting/assert"
	"github.com/speedland/go/x/xtime"
)

func TestPrinter(t *testing.T) {
	a := assert.New(t)
	s := printerFunc(func(p *Printer) {
		p.Println("foo")
	}).String()
	a.EqStr("foo\n", s)
}

func ExamplePrinter() {
	formatter := NewTextFormatter(
		`{{formattimestamp .}} [{{.Level}}] {{.Data}}`,
	)
	loggerInfo := New(
		LevelFilter(LevelInfo).Pipe(
			NewIOSinkWithFormatter(os.Stdout, formatter),
		),
	)
	loggerDebug := New(
		LevelFilter(LevelDebug).Pipe(
			NewIOSinkWithFormatter(os.Stdout, formatter),
		),
	)
	var f = func(fmt *Printer) {
		// This is not be executed.
		fmt.Printf("foo")
	}
	xtime.RunAt(
		time.Date(2016, 1, 1, 12, 10, 25, 0, time.UTC),
		func() {
			loggerInfo.Debug(f)
			loggerDebug.Debug(f)
		},
	)
	// Output:
	// 2016-01-01T12:10:25Z [debug] foo
}
