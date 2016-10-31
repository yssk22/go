package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"regexp"
	"text/template"

	"github.com/speedland/go/x/xstrings"
)

var tagRegexp = regexp.MustCompile(`([a-z0-9A-Z]+):"([^"]+)"`)

const (
	tagNameDefault         = "default"
	tagNameEnt             = "ent"
	tagValueID             = "id"
	tagValueResetIfMissing = "resetifmissing"
	tagValueForm           = "form"
	tagValueTimestamp      = "timestamp"
)

// Generator implements generator.Generator
type Generator struct {
	Package        string
	Type           string
	Fields         []*Field
	Dependencies   map[string]string
	IDField        string
	TimestampField string
}

func NewGenerator(typeName string) *Generator {
	g := &Generator{
		Type:         typeName,
		Dependencies: make(map[string]string),
	}
	g.AddDependency("fmt")
	g.AddDependencyAs("github.com/speedland/go/gae/datastore", "helper")
	g.AddDependency("github.com/speedland/go/gae/datastore/ent")
	g.AddDependency("github.com/speedland/go/gae/memcache")
	g.AddDependency("github.com/speedland/go/lazy")
	g.AddDependency("github.com/speedland/go/x/xlog")
	g.AddDependency("golang.org/x/net/context")
	g.AddDependency("google.golang.org/appengine/datastore")
	return g
}

func (g *Generator) AddDependency(pkg string) {
	g.Dependencies[pkg] = ""
}

func (g *Generator) AddDependencyAs(pkg string, as string) {
	g.Dependencies[pkg] = as
}

// Inspect implements Generator#Inspect
func (g *Generator) Inspect(node ast.Node) bool {
	if g.Fields != nil && g.Package != "" {
		// no more inspection needed.
		return true
	}
	switch node.(type) {
	case *ast.File:
		g.Package = node.(*ast.File).Name.Name
		return true
	case *ast.GenDecl:
		decl := node.(*ast.GenDecl)
		return g.inspecgGenDecl(decl)
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
func (g *Generator) GenSource(w io.Writer) error {
	if g.Fields == nil {
		return fmt.Errorf("no struct to be generated")
	}
	t := template.Must(template.New("template").Funcs(templateHelper).Parse(codeTemplate))
	return t.Execute(w, g)
}

func (g *Generator) inspecgGenDecl(decl *ast.GenDecl) bool {
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
		if t.Name.Name != g.Type {
			continue
		}
		g.Fields = make([]*Field, structType.Fields.NumFields(), structType.Fields.NumFields())
		for i := range g.Fields {
			g.Fields[i] = newField(
				g,
				structType.Fields.List[i],
			)
		}
		return true
	}
	return true
}

// Field is a type field
type Field struct {
	g *Generator
	*ast.Field
	Default        string
	IsID           bool
	IsTimestamp    bool
	ResetIfMissing bool
	AllowForm      bool
}

func (f *Field) FieldName() string {
	return f.Field.Names[0].Name
}

func (f *Field) GetDefaultExpr(v string) string {
	switch f.Field.Type.(type) {
	case *ast.Ident:
		i := f.Field.Type.(*ast.Ident)
		return builtinDefaultValue(i.Name, v)
	case *ast.SelectorExpr:
		s := f.Field.Type.(*ast.SelectorExpr)
		i := fmt.Sprintf("%s.%s", s.X, s.Sel)
		if fun, ok := defaultValueGen[i]; ok {
			dependency, expression := fun(v)
			if dependency != "" {
				f.g.AddDependency(dependency)
			}
			return expression
		}
		panic(fmt.Errorf("unsupported type: %s (%s)", i, f.FieldName()))
	default:
		panic(fmt.Errorf("unsupported expression: %s", f.FieldName()))
	}
}

func newField(g *Generator, f *ast.Field) *Field {
	field := &Field{
		g:     g,
		Field: f,
	}
	if tags := tagRegexp.FindAllStringSubmatch(f.Tag.Value, -1); tags != nil {
		for _, tag := range tags {
			tagName := tag[1]
			tagValue := tag[2]
			switch tagName {
			case tagNameDefault:
				field.Default = field.GetDefaultExpr(tagValue)
				break
			case tagNameEnt:
				values := xstrings.SplitAndTrim(tagValue, ",")
				field.IsID = hasTagValue(tagValueID, values)
				field.AllowForm = hasTagValue(tagValueForm, values)
				field.IsTimestamp = hasTagValue(tagValueTimestamp, values)
			default:
				break
			}
		}
	}
	if field.IsID {
		g.IDField = field.FieldName()
	}
	if field.IsTimestamp {
		g.TimestampField = field.FieldName()
	}
	return field
}

func hasTagValue(v string, values []string) bool {
	for _, vv := range values {
		if v == vv {
			return true
		}
	}
	return false
}
