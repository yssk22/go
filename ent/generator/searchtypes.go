package generator

type searchTypeDef struct {
	Type      string
	Converter string
}

// Man [type] => search type string
var searchTypes = map[string]*searchTypeDef{
	"bool": &searchTypeDef{
		Type:      "google.golang.org/appengine/search.Atom",
		Converter: "github.com/yssk22/go/ent.BoolToAtom",
	},
	"string": &searchTypeDef{
		Type: "string",
	},
	"[]byte": &searchTypeDef{
		Type:      "google.golang.org/appengine/search.HTML",
		Converter: "github.com/yssk22/go/ent.BytesToHTML",
	},
	"int": &searchTypeDef{
		Type:      "float64",
		Converter: "float64",
	},
	"int64": &searchTypeDef{
		Type:      "float64",
		Converter: "float64",
	},
	"float32": &searchTypeDef{
		Type:      "float64",
		Converter: "float64",
	},
	"float64": &searchTypeDef{
		Type: "float64",
	},
	"time.Time": &searchTypeDef{
		Type:      "float64",
		Converter: "github.com/yssk22/go/ent.TimeToFloat64",
	},
	"appengine.GeoPoint": &searchTypeDef{
		Type: "google.golang.org/appengine.GeoPoint",
	},
}
