package slice

import "fmt"

func ExampleSplitByLength() {
	var a = []int{0, 1, 2, 3, 4}
	var b = SplitByLength(a, 2)
	fmt.Println(b)
	// Output: [[0 1] [2 3] [4]]
}
