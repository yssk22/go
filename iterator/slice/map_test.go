package slice

import "fmt"

func ExampleMap() {
	var a = []int{0, 1, 2, 3, 4}
	b, _ := Map(a, func(i int, v int) (string, error) {
		return fmt.Sprintf("%d", v+1), nil
	})
	fmt.Println(b)
	// Output: [1 2 3 4 5]
}
