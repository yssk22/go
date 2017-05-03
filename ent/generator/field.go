package generator

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/speedland/go/x/xstrings"
)

// Field is a type field
type Field struct {
	s *Struct
	*ast.Field
	Default              string // go code expression for the default value
	Form                 string // go code expression to parse form value (var name is v and returns (value, error))
	SearchFieldTypeName  string // go code expression for SearchDoc typename
	SearchFieldConverter string // go code expression to convert a original ent field tothe serach doc field
	Parser               string // optional field to parse the string value (pkg.Func format)
	IsID                 bool
	IsTimestamp          bool
	IsForm               bool
	IsSearch             bool
	ResetIfMissing       bool
}

// FieldName returns a field name string
func (f *Field) FieldName() string {
	return f.Field.Names[0].Name
}

// FieldNameSnakeCase returns a snakecase of the field name.
func (f *Field) FieldNameSnakeCase() string {
	return xstrings.ToSnakeCase(f.FieldName())
}

// TypeName returns a type name string of the field.
func (f *Field) TypeName() string {
	var typeName string
	switch f.Field.Type.(type) {
	case *ast.Ident: // built-in
		typeName = f.Field.Type.(*ast.Ident).Name
	case *ast.SelectorExpr: // pkg.Type
		s := f.Field.Type.(*ast.SelectorExpr)
		typeName = fmt.Sprintf("%s.%s", s.X, s.Sel)
	case *ast.ArrayType: // [](type)
		at := f.Field.Type.(*ast.ArrayType)
		if i, ok := at.Elt.(*ast.Ident); ok {
			typeName = fmt.Sprintf("[]%s", i.Name)
		} else if s, ok := at.Elt.(*ast.SelectorExpr); ok {
			typeName = fmt.Sprintf("[]%s.%s", s.X, s.Sel)
		}
	}
	if typeName == "" {
		panic(fmt.Errorf("could not deletct type on %s", f.FieldName()))
	}
	return typeName
}

// GetSearchFieldType returns a type name string for search of the field.
func (f *Field) GetSearchFieldTypeName() (string, string) {
	typeName := f.TypeName()
	typeDef, ok := searchTypes[typeName]
	if !ok {
		panic(fmt.Errorf("%s (on %s) is not supported for a search field", typeName, f.FieldName()))
	}
	typePkg, typeShortExpr := parseFullPackageExpr(typeDef.Type)
	if typeShortExpr == "" {
		panic(fmt.Errorf("[BUG] invalid type definition for %s (on %s)", typeName, f.FieldName()))
	}
	convPkg, convShortExpr := parseFullPackageExpr(typeDef.Converter)
	// it's ok for convShortExpr to be "" since it means "no conversion needed".
	if typePkg != "" {
		f.s.AddDependency(typePkg)
	}
	if convPkg != "" {
		f.s.AddDependency(convPkg)
	}
	return typeShortExpr, convShortExpr
}

// GetDefaultExpr returns the default value expression for the field.
func (f *Field) GetDefaultExpr(v string) string {
	typeName := f.TypeName()
	genf, ok := defaultValueGen[typeName]
	if !ok {
		panic(fmt.Errorf("unsupported type %q on %s", typeName, f.FieldName()))
	}
	dep, expr := genf(v)
	if dep != "" {
		f.s.AddDependency(dep)
	}
	return expr
}

// GetFormExpr returns a form field expression of the field.
func (f *Field) GetFormExpr() string {
	typeName := f.TypeName()
	if typeName == "string" {
		return "v.(string)"
	}
	parser := f.Parser
	if parser == "" {
		if parser = buildInParsers[typeName]; parser == "" {
			panic(fmt.Errorf("missing parser tag for non-builtin types (%q on %s)", typeName, f.FieldName()))
		}
	}
	var fun string
	lastDot := strings.LastIndex(parser, ".")
	if lastDot == len(parser)-1 {
		panic(fmt.Errorf("parser tag must be `pkg.Func` format: %v, (%q on %s)", parser, typeName, f.FieldName()))
	} else if lastDot == -1 {
		fun = parser
	} else {
		// pkg.Func(v)
		pkg := parser[0:lastDot]
		pkgParts := strings.Split(pkg, "/")
		fun = parser[lastDot+1:]
		fun = fmt.Sprintf("%s.%s", pkgParts[len(pkgParts)-1], fun)
		f.s.AddDependency(pkg)
	}
	return fmt.Sprintf("%s(v.(string))", fun)
}

func hasTagValue(v string, values []string) bool {
	for _, vv := range values {
		if v == vv {
			return true
		}
	}
	return false
}

func parseFullPackageExpr(expr string) (string, string) {
	var shortExpr string
	var pkg string
	lastDot := strings.LastIndex(expr, ".")
	if lastDot == len(expr)-1 {
		return "", ""
	} else if lastDot == -1 {
		shortExpr = expr
	} else {
		// pkg.Type style
		pkg = expr[0:lastDot]
		pkgParts := strings.Split(pkg, "/")
		shortExpr = expr[lastDot+1:]
		shortExpr = fmt.Sprintf("%s.%s", pkgParts[len(pkgParts)-1], shortExpr)
	}
	return pkg, shortExpr
}
