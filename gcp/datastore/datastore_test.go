package datastore

import (
	"log"
	"os"
	"testing"
)

var testEnv *TestEnv

func TestMain(m *testing.M) {
	testEnv = MustNewTestEnv()
	var state int
	defer func() {
		err := testEnv.Shutdown()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(state)
	}()
	state = m.Run()
}

type Example struct {
	ID string
}
