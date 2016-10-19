package keyvalue

import "fmt"

func ExampleList() {
	m1 := Map{
		"Foo": "1",
	}
	m2 := Map{
		"Foo": "3",
		"Bar": "2",
	}
	list := NewList(m1, m2)
	fmt.Println(list.GetStringOr("Foo", "1"))
	fmt.Println(list.GetStringOr("Bar", "2"))
	fmt.Println(list.GetStringOr("Hoge", "-"))
	// Output:
	// 1
	// 2
	// -
}
