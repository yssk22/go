package config

import (
	"fmt"
	"os"
)

func ExampleEnvVar_Get() {
	os.Setenv("FOO", "envfoo")
	os.Setenv("FOO_BAR", "envfoobar")
	fmt.Println(EnvVar.GetStringOr("Foo", "invali"))
	fmt.Println(EnvVar.GetStringOr("foo", "invali"))
	fmt.Println(EnvVar.GetStringOr("foo.bar", "invali"))
	fmt.Println(EnvVar.GetStringOr("foo-bar", "invali"))
	fmt.Println(EnvVar.GetStringOr("foo bar", "invalid"))
	fmt.Println(EnvVar.GetStringOr("foo日本語bar", "invalid"))
	// Output:
	// envfoo
	// envfoo
	// envfoobar
	// envfoobar
	// envfoobar
	// invalid
}
