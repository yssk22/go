package cache

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	state := m.Run()
	os.Exit(state)
}

type Example struct {
	ID string
}
