package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/yssk22/go/x/xerrors"

	"github.com/yssk22/go/number"
)

// FileInfo is a file info that generator is analysing
type FileInfo struct {
	Path   string
	Ast    *ast.File
	Source []byte
}

// Generator is an interface to implement generator command
type Generator interface {
	Inspect(ast.Node, *FileInfo) bool
	GenSource(io.Writer) error
}

// Run runs a generator on the package directory %s and returns the formatted source
func Run(dir string, g Generator) ([]byte, error) {
	_, files, err := parsePackage(dir)
	if err != nil {
		absPath, _ := filepath.Abs(dir)
		return nil, xerrors.Wrap(err, "failed to parse package %q", absPath)
	}
	for _, f := range files {
		ast.Inspect(f.Ast, func(n ast.Node) bool {
			return g.Inspect(n, f)
		})
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

func parsePackage(dir string) (*build.Package, []*FileInfo, error) {
	pkg, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		return nil, nil, err
	}
	var files []*FileInfo
	fs := token.NewFileSet()
	for _, goFile := range pkg.GoFiles {
		path := filepath.Join(dir, goFile)
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("io error at %s: %v", goFile, err)
		}
		parsedGoFile, err := parser.ParseFile(fs, path, src, 0)
		if err != nil {
			return nil, nil, fmt.Errorf("parse error at %s: %v", goFile, err)
		}
		files = append(files, &FileInfo{
			Source: src,
			Ast:    parsedGoFile,
		})
	}
	return pkg, files, nil
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
