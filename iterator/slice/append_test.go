package slice

import "fmt"

func ExampleAppendIfMissing() {
	var a = []int{0, 1, 2, 3, 4}
	var b = AppendIfMissing(a, 0)
	var c = AppendIfMissing(a, 5)
	fmt.Println(b)
	fmt.Println(c)
	// Output:
	// [0 1 2 3 4]
	// [0 1 2 3 4 5]
}

type ComparableExample struct {
	a string
}

func (ce *ComparableExample) String() string {
	return ce.a
}

func (ce *ComparableExample) Compare(v interface{}) (int, error) {
	if ce1, ok := v.(*ComparableExample); ok {
		if ce.a == ce1.a {
			return 0, nil
		}
		if ce.a < ce1.a {
			return -1, nil
		}
		return 1, nil
	}
	return 0, ErrInvalidComparison
}

func ExampleAppendIfMissing_withPointer() {
	var a = []*ComparableExample{
		&ComparableExample{"1"},
		&ComparableExample{"2"},
	}
	var b = AppendIfMissing(a, &ComparableExample{"1"})
	var c = AppendIfMissing(a, &ComparableExample{"3"})
	fmt.Println(b)
	fmt.Println(c)
	// Output:
	// [1 2]
	// [1 2 3]
}
