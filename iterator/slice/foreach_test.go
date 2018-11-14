package slice

import "fmt"

func ExampleForEach() {
	var a = []int{0, 1, 2, 3, 4}
	ForEach(a, func(i int, v int) error {
		a[i] = v + 1
		return nil
	})
	fmt.Println(a)
	// Output: [1 2 3 4 5]
}

func ExampleForEach_error() {
	var a = []int{0, 1, 2, 3, 4}
	var err = ForEach(a, func(i int, v int) error {
		if i%2 == 0 {
			return fmt.Errorf("even error")
		}
		return nil
	})
	fmt.Println(err)
	// Output: even error (and 2 other errors)
}
