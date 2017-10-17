package xfmt

import (
	"fmt"
)

func ExamplePaddingRight() {
	fmt.Println("+", PaddingRight("hello", 10), "+")
	fmt.Println("+", PaddingRight("hellohelloh", 10), "+")
	// Output:
	// + hello      +
	// + hellohelloh +
}
func ExamplePaddingLeft() {
	s := "hello"
	fmt.Println("+", PaddingLeft(s, 10), "+")
	fmt.Println("+", PaddingLeft("hellohelloh", 10), "+")
	// Output:
	// +      hello +
	// + hellohelloh +
}
