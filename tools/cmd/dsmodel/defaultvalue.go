package main

import "fmt"

func builtinDefaultValue(typeName string, val string) string {
	switch typeName {
	case "string":
		return fmt.Sprintf("%q", val)
	default:
		return val
	}
}

// Man [type] => (dependency, expression)
var defaultValueGen = map[string](func(string) (string, string)){
	"time.Time": func(v string) (string, string) {
		const dependency = "github.com/speedland/go/x/xtime"
		if v == "$now" {
			return dependency, "xtime.Now()"
		} else if v == "$today" {
			return dependency, "xtime.Today()"
		}
		return dependency, fmt.Sprintf("xtime.ParseDateTime(%q)", v)
	},
}
