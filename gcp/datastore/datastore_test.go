package datastore

import (
	"log"
	"os"
	"testing"
)

var testEnv *TestEnv

func TestMain(m *testing.M) {
	testEnv = MustNewTestEnv()
	state := m.Run()
	err := testEnv.Shutdown()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(state)
}

type Example struct {
	ID string
}
