package slice

import "fmt"

func ExampleToMap() {
	var a = []int{0, 1, 2, 3, 4}
	var m = ToMap(a, func(i int, v int) string {
		return fmt.Sprintf("%d", v)
	})
	fmt.Println(m.(map[string]int))
	// Output: map[0:0 1:1 2:2 3:3 4:4]
}
