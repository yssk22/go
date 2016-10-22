package web

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/speedland/go/iterator/slice"
	"github.com/speedland/go/keyvalue"
)

// PathPattern is a struct to support path parameters and matches incoming request paths.
type PathPattern struct {
	source     string
	compiled   *regexp.Regexp
	paramNames []string
}

var rePathParamSomething = regexp.MustCompile(`:[a-zA-Z0-9_]+`)
var rePathParamAnything = regexp.MustCompile(`\*[a-zA-Z0-9_]*`)
var reRegexpSpecialChars = regexp.MustCompile(`\.|\[|\{|\+|\\/?`)

// MustCompilePathPattern is like CompilePathPattern but panics if an error occurrs.
func MustCompilePathPattern(pattern string) *PathPattern {
	p, err := CompilePathPattern(pattern)
	if err != nil {
		panic(err)
	}
	return p
}

// CompilePathPattern compiles the path pattern string to *PathPattern
// A path parameter name must be [a-zA-Z0-9_]+ with : and * prefix to define the matching storategy.
//
//    - /:something/ is a pattern to match something (except '/') on the path and capture the parameter value as 'something'.
//    - /*anything/ is a pattern to match anything (including '/') on the path and capture the parameter value as 'anything'
//    - /*/ is a pattern to match anything and no parameter capturing.
//
func CompilePathPattern(pattern string) (*PathPattern, error) {
	const slash = byte('/')
	const invalidPathPattern = "Routing patttern must start with '/', but got '%s'"
	if pattern[0] != slash {
		return nil, fmt.Errorf(invalidPathPattern, pattern)
	}
	var source = pattern
	// replace regexp special chars with the escaped one.
	// e.g: /path/ -> \\/path\\/
	pattern = reRegexpSpecialChars.Copy().ReplaceAllStringFunc(pattern, func(name string) string {
		return "\\" + name
	})
	// :name -> (?P<name>[^/]+)
	pattern = rePathParamSomething.Copy().ReplaceAllStringFunc(pattern, func(name string) string {
		return fmt.Sprintf("(?P<%s>[^/]+)", name[1:])
	})
	// *name -> (?P<name>.*) or * -> (.*)
	pattern = rePathParamAnything.Copy().ReplaceAllStringFunc(pattern, func(name string) string {
		if len(name) > 1 {
			return fmt.Sprintf("(?P<%s>.*)", name[1:])
		}
		return fmt.Sprintf("(.*)")
	})
	compiled, err := regexp.Compile(
		strings.Join([]string{"^", pattern, "$"}, ""),
	)
	if err != nil {
		return nil, err
	}
	names := slice.Filter(compiled.SubexpNames(), func(_ int, s string) bool {
		return s == ""
	}).([]string)

	return &PathPattern{
		source:     source,
		compiled:   compiled,
		paramNames: names,
	}, nil
}

// Match execute the matching with the given path and return the parameter values or nil
func (pattern *PathPattern) Match(path string) (*keyvalue.GetProxy, bool) {
	var matched = pattern.compiled.FindStringSubmatch(path)
	if matched == nil {
		return nil, false
	}
	var m = keyvalue.NewMap()
	var names = pattern.paramNames
	if len(names) != len(matched)-1 {
		return nil, false
	}
	for i := 1; i < len(matched); i++ {
		val := matched[i]
		// GAE server pass url encoded values to programs and clients should pass double-encoded values
		//
		// For example, the client should path /path%252Fto%252Ffoo.json
		// if they want handle /path/to/foo.json as /:param.json (set param = "path/to/foo"),
		v, err := url.QueryUnescape(val)
		if err != nil {
			return nil, false
		}
		v, err = url.QueryUnescape(v)
		if err != nil {
			return nil, false
		}
		m[names[i-1]] = v
	}
	return keyvalue.NewGetProxy(m), true
}
