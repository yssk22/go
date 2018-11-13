package datastore

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
	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xstrings"
)

const (
	signature = "datastore"

	commandParamKind = "kind"
)

var annotation = generator.NewAnnotation(
	"datastore",
)

// Generator is a generator for datastore types
type Generator struct {
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
func (g *Generator) Run(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*generator.Result, error) {
	dep := generator.NewDependency()
	dep.Add("context")
	dep.Add("google.golang.org/appengine")
	dep.Add("google.golang.org/appengine/datastore")
	dep.Add("github.com/yssk22/go/gae/memcache")
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
	if k, ok := n.Params[commandParamKind]; ok {
		spec.KindName = k
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
		fieldSpec, err := b.parseField(pkg, f, generator.ParseTag(st.Tag(i)))
		if err != nil {
			return nil, n.GenError(xerrors.Wrap(err, "could not parse the field %q", f.Name()), nil)
		}
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
	}
	if spec.KeyField == "" {
		return nil, n.GenError(fmt.Errorf("struct %s doesn't have the key field - use ent:\"key\" tag to fix", spec.StructName), nil)
	}
	if spec.TimestampField == "" {
		return nil, n.GenError(fmt.Errorf("struct %s doesn't have the timestamp field - use ent:\"timestamp\" tag to fix", spec.StructName), nil)
	}
	return &spec, nil
}

func (b *bindings) parseField(pkg *generator.PackageInfo, field *types.Var, tags keyvalue.Getter) (*FieldSpec, error) {
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
	return &f, nil
}

const (
	fieldTagName = "ent"

	fieldTagValueID        = "id"
	fieldTagValueKey       = "key"
	fieldTagValueTimestamp = "timestamp"
	fieldTagValueSearch    = "search"
)
