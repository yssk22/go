package iterator

import "fmt"

func ExampleSplitSliceByLength() {
	var a = []int{0, 1, 2, 3, 4}
	var b = SplitSliceByLength(a, 2)
	fmt.Println(b)
	// Output: [[0 1] [2 3] [4]]
}
