package xlogtest

import (
	"bytes"
	"testing"

	"github.com/yssk22/go/x/xlog"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestLogger_Name(t *testing.T) {
	a := assert.New(t)
	var buff bytes.Buffer
	xlog.New(xlog.NewIOSinkWithFormatter(
		&buff,
		xlog.NewTextFormatter(
			`{{.LoggerKey}} {{.Data}}`,
		),
	)).Infof("FOO")
	a.EqStr("github.com/yssk22/go/x/xlog/xlogtest FOO\n", buff.String())
}
