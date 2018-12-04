package flow

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xerrors"
)

const (
	signature = "flow"
)

var annotation = generator.NewAnnotation(
	"flow",
)

// Generator is a generator for Flow types
type Generator struct {
	Options *Options
}

// GetAnnotation implements generator.Generator#GetAnnotation
func (*Generator) GetAnnotation() *generator.Annotation {
	return annotation
}

// GetFormatter implements generator.Generator#GetFormatter
func (*Generator) GetFormatter() generator.Formatter {
	return generator.JavaScriptFormatter
}

// NewGenerator returns a new instance of Generator
func NewGenerator(opts *Options) *Generator {
	return &Generator{
		Options: opts,
	}
}

// Run implementes generator.Generator#Run
func (g *Generator) Run(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*generator.Result, error) {
	dep := generator.NewDependency()
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
	t := template.Must(template.New("template").Funcs(templateHelper).Parse(templateFile))
	if err = t.Execute(&buff, b); err != nil {
		return nil, xerrors.Wrap(err, "failed to run a template")
	}
	result := []*generator.Result{
		{
			Filename: "GoTypes.js",
			Source:   buff.String(),
		},
	}
	return result, nil
}

func (b *bindings) collectSpecs(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) error {
	var specs []*Spec
	var errors []error
	for _, n := range nodes {
		spec, err := b.parseAnnotatedNode(pkg, n)
		if err != nil {
			errors = append(errors, err)
		} else {
			specs = append(specs, spec)
		}
	}
	if len(errors) > 0 {
		return xerrors.MultiError(errors)
	}

	// sort
	sort.Slice(specs, func(i, j int) bool {
		a, b := specs[i], specs[j]
		return strings.Compare(string(a.TypeName), string(b.TypeName)) < 0
	})
	b.Specs = specs
	return nil
}

func (b *bindings) parseAnnotatedNode(pkg *generator.PackageInfo, n *generator.AnnotatedNode) (*Spec, error) {
	var spec Spec
	node, ok := n.Node.(*ast.GenDecl)
	if !ok {
		return nil, n.GenError(fmt.Errorf("@flow is used on non type declaration"), nil)
	}
	if node.Tok != token.TYPE {
		return nil, n.GenError(fmt.Errorf("@flow is used on non type declaration"), nil)
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
			if !f.Exported() {
				continue
			}
			ft, err := b.getFlowType(f.Type())
			if err != nil {
				return nil, n.GenError(
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
	case *types.Basic:
		o, err := b.getFlowTypeFromBasic(ut)
		if err != nil {
			return nil, n.GenError(
				xerrors.Wrap(err, "cannot resolve flowtype: %s", ut.Name()),
				nil,
			)
		}
		spec.FlowType = o
	default:
		return nil, fmt.Errorf("unsupported type: %s", reflect.TypeOf(ut))
	}
	return &spec, nil
}

func (b *bindings) getFlowType(t types.Type) (FlowType, error) {
	switch tt := t.(type) {
	case *types.Basic:
		return b.getFlowTypeFromBasic(tt)
	case *types.Pointer:
		elem, err := b.getFlowType(tt.Elem())
		if err != nil {
			return nil, err
		}
		return &FlowTypeMaybe{
			Elem: elem,
		}, nil
	case *types.Struct:
	case *types.Named:
		return b.getFlowTypeFromNamed(tt)
	}
	return nil, fmt.Errorf("unsupported")
}

func (b *bindings) getFlowTypeFromBasic(t *types.Basic) (FlowType, error) {
	switch t.Kind() {
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

func (b *bindings) getFlowTypeFromNamed(n *types.Named) (FlowType, error) {
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
		importName := b.Dependency.Add("types")
		return &FlowTypeNamed{
			Name:       n.Obj().Name(),
			ImportPath: "types",
			ImportName: importName,
		}, nil
	}
	return nil, fmt.Errorf("unsuppored named type: %s", n)
}
