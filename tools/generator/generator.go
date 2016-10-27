package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"path/filepath"
)

type Generator interface {
	Inspect(ast.Node) bool
	GenSource(io.Writer) error
}

// Run runs a generator on the package directory %s and returns the formatted source
func Run(dir string, g Generator) ([]byte, error) {
	_, files, err := parsePackage(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		ast.Inspect(f, g.Inspect)
	}
	var buff bytes.Buffer
	g.GenSource(&buff)
	formatted, err := format.Source(buff.Bytes())
	if err != nil {
		return buff.Bytes(), fmt.Errorf("generated source cannot not be compiled: %s", err)
	}
	return formatted, nil
}

func parsePackage(dir string) (*build.Package, []*ast.File, error) {
	pkg, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		return nil, nil, err
	}
	var astFiles []*ast.File
	fs := token.NewFileSet()
	for _, goFile := range pkg.GoFiles {
		parsedGoFile, err := parser.ParseFile(fs, filepath.Join(dir, goFile), nil, 0)
		if err != nil {
			return nil, nil, fmt.Errorf("parse error at %s: %v", goFile, err)
		}
		astFiles = append(astFiles, parsedGoFile)
	}
	// check the package is valid.
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{Defs: make(map[*ast.Ident]types.Object)}
	_, err = config.Check(dir, fs, astFiles, info)
	if err != nil {
		return nil, nil, err
	}
	return pkg, astFiles, nil
}
