package xstrings

import "fmt"

func ExampleSplitAndTrim() {
	fmt.Println(SplitAndTrim("a, b, c ,d,,f,", ","))
	// Output:
	// [a b c d f]
}

func ExampleSplitAndTrimAsMap() {
	m := SplitAndTrimAsMap("a, b, c ,d,,f,", ",")
	_, ok := m["a"]
	_, not := m["z"]
	fmt.Println(ok, not)
	// Output:
	// true false
}

func ExampleToSnakeCase() {
	fmt.Println(
		ToSnakeCase("A"),
		ToSnakeCase("FooBar"),
		ToSnakeCase("APIDoc"),
		ToSnakeCase("OK"),
		ToSnakeCase("IsOK"),
		ToSnakeCase("camelCase"),
	)
	// Output:
	// a foo_bar api_doc ok is_ok camel_case
}
