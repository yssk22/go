package generator

// Man [type] => func() (dependency, expression)
var formValueGen = map[string](func() (string, string)){
	"bool": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseBool(v.(string))"
	},
	"string": func() (string, string) {
		return "", "v.(string)"
	},
	"[]string": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseStringList(v.(string))"
	},
	"[]byte": func() (string, string) {
		return "", "[]byte(v.(string))"
	},
	"int": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseInt(v.(string))"
	},
	"int64": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseInt64(v.(string))"
	},
	"float32": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseFloat32(v.(string))"
	},
	"float64": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseFloat64(v.(string))"
	},
	"time.Time": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseTime(v.(string))"
	},
	"time.Duration": func() (string, string) {
		return "github.com/speedland/go/ent", "ent.ParseDuration(v.(string))"
	},
}
