package slice

import "fmt"

func ExampleFilter() {
	var a = []int{0, 1, 2, 3, 4}
	var b = Filter(a, func(i int, v int) bool {
		return v == 2
	}).([]int)
	fmt.Println(b)
	// Output: [0 1 3 4]
}
