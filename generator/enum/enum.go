package enum

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"sort"
	"strings"
	"text/template"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xstrings"
)

var annotation = generator.NewAnnotation(
	"enum",
)

const (
	signature = "enum"
)

// Generator is a generator for Enum sources
type Generator struct {
	Package      string            // package name
	Dependencies map[string]string // imports (full package path => imported name)
	Specs        []*Spec
}

// GetAnnotation implements generator.Generator#GetAnnotation
func (*Generator) GetAnnotation() *generator.Annotation {
	return annotation
}

// GetFormatter implements generator.Generator#GetFormatter
func (*Generator) GetFormatter() generator.Formatter {
	return generator.GoFormatter
}

// NewGenerator returns a new instance of Generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Run implementes generator.Generator#Run
func (enum *Generator) Run(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*generator.Result, error) {
	dep := generator.NewDependency()
	dep.Add("encoding/json")
	dep.Add("fmt")
	b := &bindings{
		Package:    pkg.Name,
		Dependency: dep,
	}
	err := b.collectSpecs(pkg, nodes)
	if err != nil {
		return nil, err
	}
	if len(b.Specs) == 0 {
		return nil, nil
	}
	var buff bytes.Buffer
	t := template.Must(template.New("template").Parse(templateFile))
	if err = t.Execute(&buff, b); err != nil {
		return nil, xerrors.Wrap(err, "failed to run a template")
	}
	result := []*generator.Result{
		{
			Filename: "generated_enums.go",
			Source:   buff.String(),
		},
	}
	return result, nil
}

func (b *bindings) collectSpecs(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) error {
	specs, err := b.collectEnumDecls(pkg, nodes)
	if err != nil {
		return err
	}
	maps := b.collectConstDelcs(pkg)
	for _, spec := range specs {
		spec.Values, err = filterValues(spec.EnumName, maps)
		if err != nil {
			return err
		}
	}
	sort.Slice(specs, func(i, j int) bool {
		a, b := specs[i], specs[j]
		return strings.Compare(a.EnumName, b.EnumName) < 0
	})
	b.Specs = specs
	return nil
}

func (b *bindings) collectEnumDecls(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*Spec, error) {
	var specs []*Spec
	for _, n := range nodes {
		node, ok := n.Node.(*ast.GenDecl)
		if !ok {
			return nil, fmt.Errorf("@enum not a decralation")
		}
		if node.Tok != token.TYPE {
			return nil, fmt.Errorf("@enum not a non type decration")
		}
		typeSpec := node.Specs[0].(*ast.TypeSpec)
		spec := &Spec{
			EnumName: typeSpec.Name.Name,
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func (b *bindings) collectConstDelcs(pkg *generator.PackageInfo) map[string][]types.Object {
	decls := make(map[string][]types.Object)
	for _, f := range pkg.Files {
		ast.Inspect(f.Ast, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				return true
			}
			for _, spec := range decl.Specs {
				vspec := spec.(*ast.ValueSpec)
				for _, n := range vspec.Names {
					typeDef := pkg.TypeInfo.Defs[n]
					typeName := typeDef.Type().String()
					if _, ok := decls[typeName]; !ok {
						decls[typeName] = []types.Object{}
					}
					decls[typeName] = append(decls[typeName], typeDef)
				}
			}
			return false
		})
	}
	return decls
}

func filterValues(enumName string, maps map[string][]types.Object) ([]Value, error) {
	constants, ok := maps[".."+enumName]
	if !ok {
		return nil, fmt.Errorf("@enum no constant is defined for %s", enumName)
	}
	var values []Value
	for _, c := range constants {
		name := c.Name()
		if !strings.HasPrefix(name, enumName) {
			return nil, fmt.Errorf("@enum each of constant name must start with %s but not (%s)", enumName, c.Name())
		}
		value := c.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
		if value.Kind() != constant.Int {
			return nil, fmt.Errorf("@enum must be an integer value for %s, but %s", enumName, value.String())
		}
		val, ok := constant.Int64Val(value)
		if !ok {
			return nil, fmt.Errorf("@enum must be an int64 value for %s, but %s", enumName, value.String())
		}
		values = append(values, Value{
			Name:     name,
			Value:    val,
			StrValue: xstrings.ToSnakeCase(strings.TrimPrefix(name, enumName)),
		})
	}
	return values, nil
}
