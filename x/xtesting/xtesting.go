package xtesting

import "flag"

var isTesting = false

func init() {
	isTesting = (flag.Lookup("test.v") != nil)
}

// IsTesting returns if the current process is executed by `go test` or not.
func IsTesting() bool {
	return isTesting
}
