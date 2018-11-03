package enum

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"sort"
	"strings"
	"text/template"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xstrings"
)

// Generator is a generator for Enum sources
// Usage:
//   @enum
//
type Generator struct {
	Package      string            // package name
	Dependencies map[string]string // imports (full package path => imported name)
	Specs        []*Spec
}

// NewGenerator returns a new instance of Generator
func NewGenerator(specs ...*Spec) *Generator {
	return &Generator{
		Dependencies: map[string]string{
			"encoding/json": "",
			"fmt":           "",
		},
		Specs: specs,
	}
}

// Run implementes generator.Generator#Run
func (enum *Generator) Run(pkg *generator.PackageInfo) ([]*generator.Result, error) {
	enum.Package = pkg.Package.Name()
	specs, err := enum.collectSpecs(pkg)
	if err != nil {
		return nil, err
	}
	enum.Specs = specs
	var buff bytes.Buffer
	t := template.Must(template.New("template").Parse(templateFile))
	if err = t.Execute(&buff, enum); err != nil {
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

func (enum *Generator) collectSpecs(pkg *generator.PackageInfo) ([]*Spec, error) {
	specs, err := enum.collectEnumDecls(pkg)
	if err != nil {
		return nil, err
	}
	maps := enum.collectConstDelcs(pkg)
	for _, spec := range specs {
		spec.Values, err = enum.filterValues(spec.EnumName, maps)
		if err != nil {
			return nil, err
		}
		log.Printf("INFO: @enum %s", spec.EnumName)
	}
	sort.Slice(specs, func(i, j int) bool {
		a, b := specs[i], specs[j]
		return strings.Compare(a.EnumName, b.EnumName) < 0
	})
	return specs, nil
}

func (enum *Generator) collectEnumDecls(pkg *generator.PackageInfo) ([]*Spec, error) {
	signatures := pkg.CollectSignatures("enum")
	var specs []*Spec
	for _, s := range signatures {
		node, ok := s.Node.(*ast.GenDecl)
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

func (enum *Generator) collectConstDelcs(pkg *generator.PackageInfo) map[string][]types.Object {
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

func (enum *Generator) filterValues(enumName string, maps map[string][]types.Object) ([]Value, error) {
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
