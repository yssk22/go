package xstrings

import "fmt"

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
