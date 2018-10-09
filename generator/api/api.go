package api

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"log"
	"text/template"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xerrors"
)

// Generator is a generator for HTTP API sources
// Usage:
//   @api path=path/to/api method=[GET|POST|PUT|DELETE]
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
			"github.com/yssk22/go/web":          "",
			"github.com/yssk22/go/web/response": "",
		},
		Specs: specs,
	}
}

// Run implementes generator.Generator#Run
func (api *Generator) Run(pkg *generator.PackageInfo) ([]*generator.Result, error) {
	api.Package = pkg.Package.Name()
	specs, err := api.collectSpecs(pkg)
	if err != nil {
		return nil, err
	}
	api.Specs = specs
	var buff bytes.Buffer
	t := template.Must(template.New("template").Funcs(templateHelper).Parse(templateFile))
	if err = t.Execute(&buff, api); err != nil {
		return nil, xerrors.Wrap(err, "failed to run a template")
	}
	result := []*generator.Result{
		{
			Filename: "__generated__apis.go",
			Source:   buff.String(),
		},
	}
	return result, nil
}

func (api *Generator) collectSpecs(pkg *generator.PackageInfo) ([]*Spec, error) {
	signatures := pkg.CollectSignatures("api")
	var specs []*Spec
	for _, s := range signatures {
		node, ok := s.Node.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcName := node.Name.Name
		path, ok := s.Params["path"]
		if !ok {
			continue
		}
		method, ok := s.Params["method"]
		if method == "" {
			method = guessMethodByFunctionName(funcName)
		}
		if method == "" {
			continue
		}
		log.Printf("INFO: @api %s %s => %s", method, path, funcName)
		specs = append(specs, &Spec{
			PathPattern: path,
			FuncName:    funcName,
			Method:      method,
		})
	}
	return specs, nil
}

var templateHelper = template.FuncMap(map[string]interface{}{
	"json": func(s *Spec) string {
		encoded, _ := json.Marshal(s)
		return string(encoded)
	},
})
