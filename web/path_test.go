package web

import "fmt"

func ExamplePathPattern_Match_Something() {
	p := MustCompilePathPattern(
		"/path/to/:something.html",
	)
	if param, ok := p.Match("/path/to/a.html"); ok {
		fmt.Printf(":something => %q", param.GetStringOr("something", ""))
	}
	// Output:
	// :something => "a"
}

func ExamplePathPattern_Match_Anything() {
	p := MustCompilePathPattern(
		"/path/to/*anything",
	)
	if param, ok := p.Match("/path/to/a.html"); ok {
		fmt.Printf(":anything => %q\n", param.GetStringOr("anything", ""))
	}
	if param, ok := p.Match("/path/to/"); ok {
		fmt.Printf(":anything => %q\n", param.GetStringOr("anything", ""))
	}
	// Output:
	// :anything => "a.html"
	// :anything => ""
}

func ExamplePathPattern_Match_URLEncoded() {
	p := MustCompilePathPattern(
		"/path/to/:something.html",
	)
	if _, ok := p.Match("/path/to/foo/bar.html"); !ok {
		fmt.Println("Not matched")
	}
	if param, ok := p.Match("/path/to/foo%2Fbar.html"); ok {
		fmt.Printf(":something => %q\n", param.GetStringOr("something", ""))
	}
	if param, ok := p.Match("/path/to/foo%252Fbar.html"); ok {
		fmt.Printf(":something => %q\n", param.GetStringOr("something", ""))
	}
	// Output:
	// Not matched
	// :something => "foo/bar"
	// :something => "foo/bar"
}
