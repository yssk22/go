package api

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"
	"text/template"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/api"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xstrings"
)

var annotation = generator.NewAnnotation(
	"api",
)

const (
	signature          = "api"
	commandParamPath   = "path"
	commandParamFormat = "format"
)

// Generator is a generator for HTTP API sources
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

// Run implementes generator.Generator#Run
func (api *Generator) Run(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*generator.Result, error) {
	dep := generator.NewDependency()
	dep.Add("github.com/yssk22/go/web")
	dep.Add("github.com/yssk22/go/web/response")
	dep.Add("github.com/yssk22/go/web/api")
	b := &bindings{
		Package: pkg.Name,
		Dependency: dep,
	}
	specs, err := b.collectSpecs(pkg, nodes)
	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return nil, nil
	}
	b.Specs = specs
	var buff bytes.Buffer
	t := template.Must(template.New("template").Funcs(templateHelper).Parse(templateFile))
	if err = t.Execute(&buff, b); err != nil {
		return nil, xerrors.Wrap(err, "failed to run a template")
	}
	result := []*generator.Result{
		{
			Filename: "generated_apis.go",
			Source:   buff.String(),
		},
	}
	return result, nil
}

func (b *bindings) collectSpecs(pkg *generator.PackageInfo, nodes []*generator.AnnotatedNode) ([]*Spec, error) {
	var specs []*Spec
	var errors []error
	for _, n := range nodes {
		spec, err := parseAnnotation(pkg, n)
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
		pathCompare := strings.Compare(a.PathPattern, b.PathPattern)
		if pathCompare != 0 {
			return pathCompare < 0
		}
		return strings.Compare(string(a.Method), string(b.Method)) < 0
	})

	// resolve dependencies
	for _, s := range specs {
		if s.StructuredParameter != nil {
			b.Dependency.Add("encoding/json")
			s.StructuredParameter.Type.ResolveAlias(b.Dependency)
		}
	}

	return specs, nil
}

func parseAnnotation(pkg *generator.PackageInfo, s *generator.AnnotatedNode) (*Spec, error) {
	var spec Spec
	node, ok := s.Node.(*ast.FuncDecl)
	if !ok {
		return nil, s.GenError(fmt.Errorf("@api is used on non func"), nil)
	}
	// check "path" parameter
	declaredParams := node.Type.Params
	pattern, err := web.CompilePathPattern(s.Params[commandParamPath])
	if err != nil {
		return nil, s.GenError(xerrors.Wrap(err, "invalid path parameter %q", s.Params[commandParamPath]), nil)
	}
	method, err := guessMethodByFunctionName(node.Name.Name)
	if err != nil {
		return nil, s.GenError(err, nil)
	}
	spec.Method = method
	spec.FuncName = node.Name.Name
	spec.PathPattern = s.Params[commandParamPath]

	// parse arguments to verify arguments' parameters mach with "path" parameter
	// arguments must be (context.Context, pathparam1, pathparam2, ..., string, (query +body struct))
	var pathParamNames = pattern.GetParamNames()
	var pathParamNamesMap = make(map[string]bool)
	for _, param := range pathParamNames {
		pathParamNamesMap[param] = true
	}
	var arguments []types.Object
	for _, paramTypeNode := range declaredParams.List {
		arguments = append(arguments, pkg.TypeInfo.Defs[paramTypeNode.Names[0]])
	}
	var hasStructuredParam = false
	if len(arguments) < (len(pathParamNames) + 1) {
		return nil, s.GenError(fmt.Errorf(
			"func %q has %d parameters, but there are only %d path parameters in the annotation",
			node.Name.Name,
			len(arguments),
			len(pathParamNames),
		), nil)
	} else if len(arguments) > (len(pathParamNames) + 1) {
		hasStructuredParam = true
	}
	if arguments[0].Type().String() != "context.Context" {
		return nil, s.GenError(fmt.Errorf(
			"func %q must have context.Context parameter in the first argument",
			node.Name.Name,
		), declaredParams.List[0])
	}
	for i := 0; i < len(pathParamNames); i++ {
		arg := arguments[i+1]
		argumentName := arg.Name()
		if arg.Type().Underlying().String() != "string" {
			return nil, s.GenError(fmt.Errorf(
				"func %q: argument %q must be string since it should be a path parameter",
				node.Name.Name,
				argumentName,
			), declaredParams.List[i+1])
		}
		if _, ok := pathParamNamesMap[argumentName]; !ok {
			return nil, s.GenError(fmt.Errorf(
				"func %q have an argument named %q, but there is no such a path parameter",
				node.Name.Name,
				argumentName,
			), nil)
		}
		spec.PathParameters = append(spec.PathParameters, argumentName)
	}
	if hasStructuredParam {
		format, _ := api.ParseRequestParameterFormat(s.Params[commandParamFormat])
		arg := arguments[len(pathParamNames)+1]
		parameter, err := getParameterParser(arg, resolveRequestParameterFormat(format, spec.Method))
		if err != nil {
			return nil, s.GenError(xerrors.Wrap(err, "could not build parameter parser"), nil)
		}
		spec.StructuredParameter = parameter
	}
	return &spec, nil
}

func getParameterParser(arg types.Object, format api.RequestParameterFormat) (*StructuredParameter, error) {
	p, ok := arg.Type().(*types.Pointer)
	if !ok {
		return nil, fmt.Errorf("%s must be a pointer of named struct but %s", arg.Name(), arg.Type().String())
	}
	n, ok := p.Elem().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("%s must be a pointer of named struct but %s", arg.Name(), p.Elem().String())
	}
	obj := n.Obj()
	pkgPath := obj.Pkg().Path()
	if pkgPath == "." {
		pkgPath = ""
	}
	var s = StructuredParameter{
		Type: &ParameterType{
			Name:    obj.Name(),
			Package: pkgPath,
		},
		Parser: api.NewParameterParser(format),
	}
	st, ok := n.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s must be a struct but %s", s.Type, n.Underlying().String())
	}
	l := st.NumFields()
	for i := 0; i < l; i++ {
		err := configureParameterParserForField(s.Parser, st.Field(i), generator.ParseTag(st.Tag(i)))
		if err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func guessMethodByFunctionName(funcName string) (requestMethod, error) {
	if strings.HasPrefix(funcName, "get") {
		return requestMethodGet, nil
	}
	if strings.HasPrefix(funcName, "update") {
		return requestMethodPut, nil
	}
	if strings.HasPrefix(funcName, "create") {
		return requestMethodPost, nil
	}
	if strings.HasPrefix(funcName, "delete") {
		return requestMethodDelete, nil
	}
	return requestMethodUnknown, fmt.Errorf("invalid function name %q to resolve the method", funcName)
}

func resolveRequestParameterFormat(format api.RequestParameterFormat, m requestMethod) api.RequestParameterFormat {
	if format.IsVaild() {
		return format
	}
	switch m {
	case requestMethodGet:
		return api.RequestParameterFormatQuery
	case requestMethodDelete:
		return api.RequestParameterFormatQuery
	case requestMethodPut:
		if format == api.RequestParameterFormatQuery || format == api.RequestParameterFormatJSON {
			return format
		}
		return api.RequestParameterFormatForm
	case requestMethodPost:
		if format == api.RequestParameterFormatQuery || format == api.RequestParameterFormatJSON {
			return format
		}
		return api.RequestParameterFormatForm
	}
	return api.RequestParameterFormatQuery
}

func configureParameterParserForField(pp *api.ParameterParser, field *types.Var, tags keyvalue.Getter) error {
	var parameterName string
	if v, err := tags.Get("json"); err == nil {
		values := strings.Split(v.(string), ",")
		parameterName = values[0]
	} else {
		parameterName = xstrings.ToSnakeCase(field.Name())
	}
	t, err := getParameterType(field.Type())
	if err != nil {
		return fmt.Errorf("field %s: couldn't resolve the type - %s", field.Name(), err)
	}
	pp.Type(parameterName, t)
	if v, err := tags.Get("validate"); err == nil {
		values := strings.Split(v.(string), ",")
		for _, v := range values {
			if v == "required" {
				pp.Required(parameterName)
			}
		}
	}

	if v, err := tags.Get("default"); err == nil {
		value, err := t.ValueOf(v.(string))
		if err != nil {
			return fmt.Errorf("field %s: invalid default value - %s", field.Name(), err)
		}
		pp.Default(parameterName, value)
	}
	return nil
}

func getParameterType(t types.Type, upperTypes ...string) (api.RequestParameterFieldType, error) {
	s := t.String()
	switch s {
	case "string":
		return api.RequestParameterFieldTypeString, nil
	case "int":
		return api.RequestParameterFieldTypeInt, nil
	case "float64":
		return api.RequestParameterFieldTypeFloat, nil
	case "time.Time":
		return api.RequestParameterFieldTypeTime, nil
	case "bool":
		return api.RequestParameterFieldTypeBool, nil
	}
	switch t.(type) {
	case *types.Array:
		return api.RequestParameterFieldTypeArray, nil
	case *types.Slice:
		return api.RequestParameterFieldTypeArray, nil
	case *types.Struct:
		return api.RequestParameterFieldTypeObject, nil
	case *types.Pointer:
		return getParameterType(t.(*types.Pointer).Elem(), append(upperTypes, t.String())...)
	}
	if t != t.Underlying() {
		return getParameterType(t.Underlying(), append(upperTypes, t.String())...)
	}
	return api.RequestParameterFieldType(0), fmt.Errorf("unresolved type %s", strings.Join(append(upperTypes, t.String()), "->"))
}
