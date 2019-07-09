package xerrors

import (
	"fmt"
	"strings"
)

func ExampleWrapAndUnwrap() {
	rootError := fmt.Errorf("root error")
	error1st := Wrap(rootError, "1st error")
	list := Unwrap(error1st)

	// don't test with source line
	firstLine := "-" + strings.Split(list[0], "-")[1]
	lastLine := list[1]
	fmt.Println(firstLine)
	fmt.Println(lastLine)
	// Output:
	// - 1st error
	// unknown - root error
}
