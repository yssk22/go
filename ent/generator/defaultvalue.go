package generator

import "fmt"

// Man [type] => (dependency, expression)
var defaultValueGen = map[string](func(string) (string, string)){
	"string": func(v string) (string, string) {
		return "", fmt.Sprintf("%q", v)
	},
	"bool": func(v string) (string, string) {
		return "", v
	},
	"[]string": func(v string) (string, string) {
		return "github.com/speedland/go/ent", fmt.Sprintf("ent.ParseStringList(%q)", v)
	},
	"[]byte": func(v string) (string, string) {
		return "", fmt.Sprintf("[]byte(%q)", v)
	},
	"int": func(v string) (string, string) {
		return "", v
	},
	"int64": func(v string) (string, string) {
		return "", v
	},
	"float32": func(v string) (string, string) {
		return "", v
	},
	"float64": func(v string) (string, string) {
		return "", v
	},
	"time.Time": func(v string) (string, string) {
		return "github.com/speedland/go/ent", fmt.Sprintf("ent.ParseTime(%q)", v)
	},
	"time.Duration": func(v string) (string, string) {
		return "github.com/speedland/go/ent", fmt.Sprintf("ent.ParseDuration(%q)", v)
	},
}
