package number

import "fmt"

func ExampleParseFloatOr() {
	fmt.Println(ParseFloat32Or("20.1", 1))
	fmt.Println(ParseFloat32Or("20.5a", 1.3))
	fmt.Println(ParseFloat32Or("a", 1.2))
	// Output:
	// 20.1
	// 1.3
	// 1.2
}
