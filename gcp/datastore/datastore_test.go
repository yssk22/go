package datastore

import (
	"log"
	"os"
	"testing"
)

var testEnvironment *TestEnviornment

func TestMain(m *testing.M) {
	testEnvironment = MustNewTestEnviornment()
	state := m.Run()
	err := testEnvironment.Shutdown()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(state)
}

type Example struct {
	ID string
}
