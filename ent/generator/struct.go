package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"text/template"

	"github.com/speedland/go/x/xstrings"
)

// Struct is a parsed result of struct to generate a code.
type Struct struct {
	Package        string
	Type           string
	Kind           string
	Instance       string
	Fields         []*Field
	Dependencies   map[string]string
	IDField        string
	TimestampField string
	IsSearchable   bool
}

// NewStruct returns a struct for the type `typeName`
func NewStruct(typeName string, kindName string) *Struct {
	s := &Struct{
		Type:         typeName,
		Kind:         kindName,
		Instance:     xstrings.ToSnakeCase(typeName)[:1],
		Dependencies: make(map[string]string),
	}
	s.AddDependency("fmt")
	s.AddDependencyAs("github.com/speedland/go/gae/datastore", "helper")
	s.AddDependency("github.com/speedland/go/ent")
	s.AddDependency("github.com/speedland/go/gae/memcache")
	s.AddDependency("github.com/speedland/go/keyvalue")
	s.AddDependency("github.com/speedland/go/lazy")
	s.AddDependency("github.com/speedland/go/x/xlog")
	s.AddDependency("github.com/speedland/go/x/xtime")
	s.AddDependency("golang.org/x/net/context")
	s.AddDependency("google.golang.org/appengine/datastore")
	return s
}

// AddDependency adds a dependent package of the struct
func (s *Struct) AddDependency(pkg string) {
	s.Dependencies[pkg] = ""
}

// AddDependencyAs is like AddDependency with specifiyng the alias name.
func (s *Struct) AddDependencyAs(pkg string, as string) {
	s.Dependencies[pkg] = as
}

// Inspect implements Generator#Inspect
func (s *Struct) Inspect(node ast.Node) bool {
	if s.Fields != nil && s.Package != "" {
		// no more inspection needed.
		return true
	}
	switch node.(type) {
	case *ast.File:
		s.Package = node.(*ast.File).Name.Name
		return true
	case *ast.GenDecl:
		decl := node.(*ast.GenDecl)
		return s.inspectGenDecl(decl)
	default:
		return true
	}
}

var templateHelper = template.FuncMap(map[string]interface{}{
	"snakecase": func(s string) string {
		return xstrings.ToSnakeCase(s)
	},
})

// GenSource implements Generator#GenSource
func (s *Struct) GenSource(w io.Writer) error {
	if s.Fields == nil {
		return fmt.Errorf("no struct to be generated")
	}
	t := template.Must(template.New("template").Funcs(templateHelper).Parse(codeTemplate))
	return t.Execute(w, s)
}

func (s *Struct) inspectGenDecl(decl *ast.GenDecl) bool {
	if decl.Tok != token.TYPE {
		return true
	}
	for _, spec := range decl.Specs {
		t, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		structType, ok := t.Type.(*ast.StructType)
		if !ok {
			continue
		}
		if t.Name.Name != s.Type {
			continue
		}
		s.Fields = make([]*Field, structType.Fields.NumFields(), structType.Fields.NumFields())
		for i := range s.Fields {
			s.Fields[i] = s.newField(
				structType.Fields.List[i],
			)
			if s.Fields[i].IsSearch {
				s.IsSearchable = true
				s.AddDependency("google.golang.org/appengine/search")
			}
		}
		return true
	}
	return true
}

func (s *Struct) newField(f *ast.Field) *Field {
	field := &Field{
		s:     s,
		Field: f,
	}
	var defaultValue string
	if f.Tag != nil {
		if tags := tagRegexp.FindAllStringSubmatch(f.Tag.Value, -1); tags != nil {
			for _, tag := range tags {
				tagName := tag[1]
				tagValue := tag[2]
				switch tagName {
				case tagNameParser:
					field.Parser = tagValue
				case tagNameDefault:
					defaultValue = tagValue
					break
				case tagNameEnt:
					values := xstrings.SplitAndTrim(tagValue, ",")
					field.IsID = hasTagValue(tagValueID, values)
					field.IsTimestamp = hasTagValue(tagValueTimestamp, values)
					field.IsForm = hasTagValue(tagValueForm, values)
					field.ResetIfMissing = hasTagValue(tagValueResetIfMissing, values)
					field.IsSearch = hasTagValue(tagValueSearch, values)
				default:
					break
				}
			}
		}
	}
	// TODO: Refactoring the field generation here!!!!! avoid modification both on f and s.
	if field.IsID {
		s.IDField = field.FieldName()
	}
	if field.IsTimestamp {
		s.TimestampField = field.FieldName()
	}
	if field.IsSearch {
		field.SearchFieldTypeName, field.SearchFieldConverter = field.GetSearchFieldTypeName()
	}
	if field.IsForm {
		field.Form = field.GetFormExpr()
	}
	if defaultValue != "" {
		field.Default = field.GetDefaultExpr(defaultValue)
	}
	return field
}
