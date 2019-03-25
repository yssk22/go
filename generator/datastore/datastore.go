package datastore

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
	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xstrings"
)

const (
	signature = "datastore"

	commandParamKind = "kind"
)

var annotation = generator.NewAnnotationSymbol("datastore")

// Generator is a generator for datastore types
type Generator struct {
}

// GetAnnotationSymbol implements generator.Generator#GetAnnotationSymbol
func (*Generator) GetAnnotationSymbol() generator.AnnotationSymbol {
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
func (g *Generator) Run(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*generator.Result, error) {
	dep := generator.NewDependency()
	dep.Add("context")
	dep.Add("google.golang.org/appengine/datastore")
	dep.AddAs("github.com/yssk22/go/gae/datastore", "ds")
	dep.Add("github.com/yssk22/go/x/xerrors")
	dep.Add("github.com/yssk22/go/x/xtime")
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
			Filename: "generated_datastore.go",
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
		return strings.Compare(string(a.StructName), string(b.StructName)) < 0
	})
	b.Specs = specs
	return nil
}

func (b *bindings) parseAnnotatedNode(pkg *generator.PackageInfo, n *generator.AnnotatedNode) (*Spec, error) {
	var spec Spec
	node, ok := n.Node.(*ast.GenDecl)
	if !ok {
		return nil, n.GenError(fmt.Errorf("@datastore is used on non type declaration"), nil)
	}
	if node.Tok != token.TYPE {
		return nil, n.GenError(fmt.Errorf("@datastore is used on non type declaration"), nil)
	}
	typeSpec := node.Specs[0].(*ast.TypeSpec)
	t := pkg.TypeInfo.Defs[typeSpec.Name]
	spec.StructName = t.Name()
	params := n.GetParamsBy(annotation)
	if k, err := params.Get(commandParamKind); err == nil {
		spec.KindName = k.(string)
	} else {
		spec.KindName = spec.StructName
	}
	st, ok := t.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, n.GenError(fmt.Errorf("@datastore is used on non struct type declaration"), nil)
	}
	l := st.NumFields()
	for i := 0; i < l; i++ {
		f := st.Field(i)
		if !f.Exported() {
			continue
		}
		tag := generator.ParseTag(st.Tag(i))
		fieldSpec, err := b.getFieldSpec(pkg, f, tag)
		if err != nil {
			return nil, n.GenError(xerrors.Wrap(err, "could not get the field spec for %q", f.Name()), nil)
		}
		if fieldSpec != nil {
			if fieldSpec.IsKey {
				if spec.KeyField != "" {
					return nil, n.GenError(fmt.Errorf("struct %s have multiple key fields - use ent:\"key\" tag only once", spec.StructName), nil)
				}
				spec.KeyField = fieldSpec.Name
			}
			if fieldSpec.IsTimestamp {
				if spec.TimestampField != "" {
					return nil, n.GenError(fmt.Errorf("struct %s have multiple timestamp fields - use ent:\"key\" tag only once", spec.StructName), nil)
				}
				spec.TimestampField = fieldSpec.Name
			}

			spec.Fields = append(spec.Fields, fieldSpec)
			if !fieldSpec.NoIndex {
				querySpecs, err := b.getQuerySpecs(pkg, f, tag)
				if err != nil {
					return nil, n.GenError(xerrors.Wrap(err, "could not get the query specs for %q", f.Name()), nil)
				}
				spec.QuerySpecs = append(spec.QuerySpecs, querySpecs...)
			}
		}
	}
	if spec.KeyField == "" {
		return nil, n.GenError(fmt.Errorf("struct %s doesn't have the key field - use ent:\"key\" tag to fix", spec.StructName), nil)
	}
	return &spec, nil
}

func (b *bindings) getFieldSpec(pkg *generator.PackageInfo, field *types.Var, tags keyvalue.Getter) (*FieldSpec, error) {
	var f FieldSpec
	f.Name = field.Name()
	if v, err := tags.Get(fieldTagName); err == nil {
		values := xstrings.SplitAndTrim(v.(string), ",")
		for _, v := range values {
			switch v {
			case fieldTagValueKey:
				f.IsKey = true
			case fieldTagValueID:
				f.IsKey = true
			case fieldTagValueSearch:
				f.IsSearch = true
			case fieldTagValueTimestamp:
				f.IsTimestamp = true
			}
		}
	}
	if v, err := tags.Get(datastoreTagName); err == nil {
		var fieldName string
		var fieldOption string
		values := strings.Split(v.(string), ",")
		fieldName = strings.TrimSpace(values[0])
		if len(values) > 1 {
			fieldOption = strings.TrimSpace(values[1])
		}
		if fieldName == "-" {
			return nil, nil
		}
		if fieldName != "" {
			if fieldName == "noindex" {
				return nil, fmt.Errorf("%s has tagged with datastore but name is spcified with noindex. You would probably want to tag \",noindex\"", f.Name)
			}
			f.Name = fieldName
		}
		if fieldOption == "noindex" {
			f.NoIndex = true
		}
	}
	return &f, nil
}

func (b *bindings) getQuerySpecs(pkg *generator.PackageInfo, field *types.Var, tags keyvalue.Getter) ([]*QuerySpec, error) {
	name := field.Name()
	t := field.Type()
	return b.getQuerySpecsRec(pkg, name, name, t)
}

func (b *bindings) getQuerySpecsRec(pkg *generator.PackageInfo, name string, propertyName string, t types.Type) ([]*QuerySpec, error) {
	var specs []*QuerySpec
	switch tt := t.(type) {
	case *types.Basic:
		specs = append(specs, &QuerySpec{
			Name:         name,
			PropertyName: propertyName,
			Type:         tt.Name(),
		})
		return specs, nil
	case *types.Pointer:
		return b.getQuerySpecsRec(pkg, name, name, tt.Elem())
	case *types.Struct:
		numFields := tt.NumFields()
		for i := 0; i < numFields; i++ {
			f := tt.Field(i)
			if !f.Exported() {
				continue
			}
			underlyingSpecs, err := b.getQuerySpecsRec(
				pkg,
				fmt.Sprintf("%s%s", name, f.Name()),
				fmt.Sprintf("%s.%s", propertyName, f.Name()),
				f.Type(),
			)
			if err != nil {
				return nil, err
			}
			specs = append(specs, underlyingSpecs...)
		}
		return specs, nil
	case *types.Named:
		str := tt.String()
		if str == "time.Time" {
			alias := b.Dependency.Add("time")
			specs = append(specs, &QuerySpec{
				Name:         name,
				PropertyName: propertyName,
				Type:         fmt.Sprintf("%s.Time", alias),
			})
			return specs, nil
		}
		underlying := tt.Underlying()
		if ut, ok := underlying.(*types.Struct); ok {
			return b.getQuerySpecsRec(pkg, name, propertyName, ut)
		}
		obj := tt.Obj()
		importPath := obj.Pkg().Path()
		if importPath == pkg.Package.Path() {
			specs = append(specs, &QuerySpec{
				Name:         name,
				PropertyName: propertyName,
				Type:         obj.Name(),
			})
		} else {
			alias := b.Dependency.Add(importPath)
			specs = append(specs, &QuerySpec{
				Name:         name,
				PropertyName: propertyName,
				Type:         fmt.Sprintf("%s.%s", alias, obj.Name()),
			})
		}
		return specs, nil
	case *types.Slice:
		str := tt.String()
		if str != "[]byte" {
			return b.getQuerySpecsRec(pkg, name, propertyName, tt.Elem())
		}
		return specs, nil
	default:
		return nil, fmt.Errorf("unsupported field %s (node type: %s)", name, reflect.ValueOf(tt))
	}
}

const (
	fieldTagName = "ent"

	fieldTagValueID        = "id"
	fieldTagValueKey       = "key"
	fieldTagValueTimestamp = "timestamp"
	fieldTagValueSearch    = "search"

	datastoreTagName = "datastore"
)
