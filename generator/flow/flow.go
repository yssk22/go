package flow

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strings"
	"text/template"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xerrors"
)

const (
	signature = "flow"
)

// Generator is a generator for Flow types
// Usage: @flow
type Generator struct {
	Package    string // package name
	Dependency *generator.Dependency
	Specs      []*Spec
}

// NewGenerator returns a new instance of Generator
func NewGenerator(specs ...*Spec) *Generator {
	dep := generator.NewDependency()
	return &Generator{
		Dependency: dep,
		Specs:      specs,
	}
}

// Run implementes generator.Generator#Run
func (g *Generator) Run(pkg *generator.PackageInfo) ([]*generator.Result, error) {
	g.Package = pkg.Package.Name()
	specs, err := g.collectSpecs(pkg)
	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return nil, nil
	}
	g.Specs = specs
	var buff bytes.Buffer
	t := template.Must(template.New("template").Funcs(templateHelper).Parse(templateFile))
	if err = t.Execute(&buff, g); err != nil {
		return nil, xerrors.Wrap(err, "failed to run a template")
	}
	result := []*generator.Result{
		{
			Filename: "ServerTypes.js",
			Source:   buff.String(),
			FileType: generator.ResultFileTypeFlow,
		},
	}
	return result, nil
}

func (g *Generator) collectSpecs(pkg *generator.PackageInfo) ([]*Spec, error) {
	signatures := pkg.CollectSignatures(signature)
	var specs []*Spec
	var errors []error
	for _, s := range signatures {
		spec, err := g.parseSignature(pkg, s)
		if err != nil {
			errors = append(errors, err)
		} else {
			specs = append(specs, spec)
		}
	}
	if len(errors) > 0 {
		return nil, xerrors.MultiError(errors)
	}

	// sort
	sort.Slice(specs, func(i, j int) bool {
		a, b := specs[i], specs[j]
		return strings.Compare(string(a.TypeName), string(b.TypeName)) < 0
	})
	return specs, nil
}

func (g *Generator) parseSignature(pkg *generator.PackageInfo, s *generator.Signature) (*Spec, error) {
	var spec Spec
	node, ok := s.Node.(*ast.GenDecl)
	if !ok {
		return nil, s.GenError(fmt.Errorf("@flow is used on non struct type definition"), nil)
	}
	if node.Tok != token.TYPE {
		return nil, s.GenError(fmt.Errorf("@flow is used on non struct type definition"), nil)
	}
	typeSpec := node.Specs[0].(*ast.TypeSpec)
	t := pkg.TypeInfo.Defs[typeSpec.Name]
	spec.TypeName = t.Name()
	switch ut := t.Type().Underlying().(type) {
	case *types.Struct:
		o := &FlowTypeObject{}
		l := ut.NumFields()
		for i := 0; i < l; i++ {
			f := ut.Field(i)
			ft, err := g.getFlowType(f.Type())
			if err != nil {
				return nil, s.GenError(
					xerrors.Wrap(err, "cannot resolve flowtype for the field %s - ", f.Name()),
					nil,
				)
			}
			o.Fields = append(o.Fields, FlowTypeObjectField{
				Name: f.Name(),
				Type: ft,
			})
		}
		spec.FlowType = o
	}
	return &spec, nil
}

func (g *Generator) getFlowType(t types.Type) (FlowType, error) {
	switch tt := t.(type) {
	case *types.Basic:
		return g.getFlowTypeFromBasic(tt)
	case *types.Pointer:
		elem, err := g.getFlowType(tt.Elem())
		if err != nil {
			return nil, err
		}
		return &FlowTypeMaybe{
			Elem: elem,
		}, nil
	case *types.Struct:
	case *types.Named:
		return g.getFlowTypeFromNamed(tt)
	}
	return nil, fmt.Errorf("unsupported")
}

func (g *Generator) getFlowTypeFromBasic(b *types.Basic) (FlowType, error) {
	switch b.Kind() {
	case types.Bool:
		return FlowTypeBool, nil
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		return FlowTypeNumber, nil
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		return FlowTypeNumber, nil
	case types.Float32, types.Float64:
		return FlowTypeNumber, nil
	case types.String:
		return FlowTypeString, nil
	}
	return nil, fmt.Errorf("unsuppored basic type: %s", b)
}

func (g *Generator) getFlowTypeFromNamed(n *types.Named) (FlowType, error) {
	s := n.String()
	switch s {
	case "time.Time":
		return FlowTypeDate, nil
	}
	importPath := n.Obj().Pkg().Path()
	if importPath == "." {
		return &FlowTypeNamed{
			Name: n.Obj().Name(),
		}, nil
	}
	switch importPath {
	case "github.com/yssk22/go/types":
		importName := g.Dependency.Add("types")
		return &FlowTypeNamed{
			Name:       n.Obj().Name(),
			ImportPath: "types",
			ImportName: importName,
		}, nil
	}
	return nil, fmt.Errorf("unsuppored named type: %s", n)
}
