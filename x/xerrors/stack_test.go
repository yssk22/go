package xerrors

import (
	"fmt"
	"strings"
)

func ExampleWrapAndUnwrap() {
	rootError := fmt.Errorf("root error")
	error1st := Wrap(rootError, "1st error")
	fmt.Println(strings.Join(Unwrap(error1st), "\n"))
	// Output:
	// /Users/yohei/sites/yssk22.dev/goss/x/xerrors/stack_test.go:10 - 1st error
	// unknown - root error
}
