package number

import "fmt"

func ExampleParseFloatOr() {
	fmt.Println(ParseFloatOr("20.1", 1))
	fmt.Println(ParseFloatOr("20.5a", 1.3))
	fmt.Println(ParseFloatOr("a", 1.2))
	// Output:
	// 20.1
	// 1.3
	// 1.2
}
