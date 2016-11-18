package generator

// Man [type] => string
var buildInParsers = map[string]string{
	"bool":          "github.com/speedland/go/ent.ParseBool",
	"[]string":      "github.com/speedland/go/ent.ParseStringList",
	"[]byte":        "[]byte",
	"int":           "github.com/speedland/go/ent.ParseInt",
	"int64":         "github.com/speedland/go/ent.ParseInt64",
	"float32":       "github.com/speedland/go/ent.ParseFloat32",
	"float64":       "github.com/speedland/go/ent.ParseFloat64",
	"time.Time":     "github.com/speedland/go/ent.ParseTime",
	"time.Duration": "github.com/speedland/go/ent.ParseDuration",
}
