package number

import "fmt"

func ExampleParseIntOr() {
	fmt.Println(ParseIntOr("20", 1))
	fmt.Println(ParseIntOr("20a", 1))
	fmt.Println(ParseIntOr("a", 1))
	// Output:
	// 20
	// 1
	// 1
}
