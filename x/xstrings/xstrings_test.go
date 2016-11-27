package xstrings

import "fmt"

func ExampleStripAndTrim() {
	fmt.Println(SplitAndTrim("a, b, c ,d,,f,", ","))
	// Output:
	// [a b c d f]
}

func ExampleToSnakeCase() {
	fmt.Println(
		ToSnakeCase("FooBar"),
		ToSnakeCase("APIDoc"),
		ToSnakeCase("OK"),
		ToSnakeCase("IsOK"),
		ToSnakeCase("camelCase"),
	)
	// Output:
	// foo_bar api_doc ok is_ok camel_case
}
