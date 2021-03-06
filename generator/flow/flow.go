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
	"github.com/yssk22/go/generator/enum"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xstrings"
)

const (
	signature = "flow"
)

var annotation = generator.NewAnnotationSymbol("flow")

// Generator is a generator for Flow types
type Generator struct {
	Options *Options
}

// GetAnnotationSymbol implements generator.Generator#AnnotationSymbol
func (*Generator) GetAnnotationSymbol() generator.AnnotationSymbol {
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
	enumValues := make(map[string]*enum.Spec)
	specs, err := enum.CollectSpecs(pkg, nodes)
	if err != nil {
		return nil, err
	}
	for _, s := range specs {
		enumValues[s.EnumName] = s
	}
	err = b.collectSpecs(pkg, nodes, enumValues)
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

func (b *bindings) collectSpecs(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode, enumSpecs map[string]*enum.Spec) error {
	var specs []*Spec
	var errors []error
	for _, n := range nodes {
		spec, err := b.parseAnnotatedNode(pkg, n, enumSpecs)
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

func (b *bindings) parseAnnotatedNode(pkg *generator.PackageInfo, n *generator.AnnotatedNode, enumSpecs map[string]*enum.Spec) (*Spec, error) {
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
	if enumSpec, ok := enumSpecs[spec.TypeName]; ok {
		spec.FlowType = &FlowTypeEnum{
			spec: enumSpec,
		}
		return &spec, nil
	}
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
			fieldName := f.Name()
			omitEmpty := false
			tags := generator.ParseTag(ut.Tag(i))
			if jsonName, err := tags.Get("json"); err == nil {
				jsonTags := xstrings.SplitAndTrim(jsonName.(string), ",")
				l := len(jsonTags)
				fieldName = jsonTags[0]
				if l == 2 {
					omitEmpty = jsonTags[1] == "omitempty"
				}
				if fieldName == "-" {
					continue
				}
			}
			o.Fields = append(o.Fields, FlowTypeObjectField{
				Name:      fieldName,
				Type:      ft,
				OmitEmpty: omitEmpty,
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
	return nil, fmt.Errorf("unsuppored basic type: %v", b)
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
	// TODO: we need to check import path and import name if named type n is defiend in the different package.
	return &FlowTypeNamed{
		Name: n.Obj().Name(),
	}, nil
}
