package gaetest

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(Run(func() int {
		return m.Run()
	}))
}
