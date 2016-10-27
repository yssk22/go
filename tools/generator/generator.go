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
	"strings"

	"github.com/speedland/go/number"
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
		return nil, &InvalidSourceError{
			Source: buff.String(),
			err:    err,
		}
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

// InvalidSourceError is an error if generated source is unable to compile.
// Use SourceWithLine() to debug the generated code.
type InvalidSourceError struct {
	Source string
	err    error
}

// SourceWithLine returns the source code with line numbers
func (e *InvalidSourceError) SourceWithLine(all bool) string {
	lines := strings.Split(e.Source, "\n")
	max := len(fmt.Sprintf("%d", len(lines)))
	format := fmt.Sprintf("%%%dd: %%s", max+1)
	errLineFormat := fmt.Sprintf("!%%%dd: %%s", max)
	n := number.ParseIntOr(strings.Split(e.err.Error(), ":")[0], 0)
	for i := range lines {
		if n != 0 && n == i+1 {
			lines[i] = fmt.Sprintf(errLineFormat, i+1, lines[i])
		} else {
			lines[i] = fmt.Sprintf(format, i+1, lines[i])
		}
	}
	if all || n == 0 {
		return strings.Join(lines, "\n")
	}
	n++
	if n < 5 {
		lines = lines[:n+4]
	} else if n < len(lines)-3 {
		lines = lines[n-6 : n+3]
	} else {
		lines = lines[n-6:]
	}
	return strings.Join(lines, "\n")
}

func (e *InvalidSourceError) Error() string {
	return e.err.Error()
}
