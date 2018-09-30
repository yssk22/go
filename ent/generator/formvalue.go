package generator

// Man [type] => string
var buildInParsers = map[string]string{
	"bool":          "github.com/yssk22/go/ent.ParseBool",
	"[]string":      "github.com/yssk22/go/ent.ParseStringList",
	"[]byte":        "[]byte",
	"int":           "github.com/yssk22/go/ent.ParseInt",
	"int64":         "github.com/yssk22/go/ent.ParseInt64",
	"float32":       "github.com/yssk22/go/ent.ParseFloat32",
	"float64":       "github.com/yssk22/go/ent.ParseFloat64",
	"time.Time":     "github.com/yssk22/go/ent.ParseTime",
	"time.Duration": "github.com/yssk22/go/ent.ParseDuration",
}
