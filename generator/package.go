package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"

	"github.com/yssk22/go/x/xerrors"
)

// PackageInfo is a package infomaiton that generator is analysing
type PackageInfo struct {
	Name     string
	Package  *types.Package
	TypeInfo *types.Info
	Files    []*FileInfo
}

func parsePackage(dir string) (*PackageInfo, error) {
	importedPackage, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		return nil, xerrors.Wrap(err, "build.Default.Import failed")
	}
	fs := token.NewFileSet()
	var parsedFiles []*ast.File
	var files []*FileInfo
	var offset = 0
	for _, goFile := range importedPackage.GoFiles {
		path := filepath.Join(dir, goFile)
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, xerrors.Wrap(err, "io error: %s", path)
		}
		parsedFile, err := parser.ParseFile(fs, path, src, parser.ParseComments)
		if err != nil {
			return nil, xerrors.Wrap(err, "parse error: %s", path)
		}
		files = append(files, &FileInfo{
			Path:       path,
			Source:     src,
			NodeOffset: offset,
			Ast:        parsedFile,
			CommentMap: ast.NewCommentMap(fs, parsedFile, parsedFile.Comments),
		})
		parsedFiles = append(parsedFiles, parsedFile)
		offset += len(src)
	}

	info := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{},
		Defs:  map[*ast.Ident]types.Object{},
		Uses:  map[*ast.Ident]types.Object{},
	}
	conf := types.Config{
		FakeImportC: true,
		Importer:    importer.Default(),
	}
	pkg, err := conf.Check(".", fs, parsedFiles, info)
	if err != nil {
		conf = types.Config{
			FakeImportC: true,
			Importer:    importer.For("source", nil),
		}
		pkg, err = conf.Check(".", fs, parsedFiles, info)
		if err != nil {
			return nil, xerrors.Wrap(err, "type check error -- you may need to run `go install ./...` at first")
		}
	}
	pkg.SetName(importedPackage.Name)
	return &PackageInfo{
		Name:     importedPackage.Name,
		Package:  pkg,
		TypeInfo: info,
		Files:    files,
	}, nil
}

// Inspect runs a ast.Inspect for files in the package
func (p *PackageInfo) Inspect(fun func(ast.Node) bool) {
	for _, f := range p.Files {
		f.Inspect(fun)
	}
}

// FileInfo is a file info that generator is analysing
type FileInfo struct {
	Path       string
	Ast        *ast.File
	Source     []byte
	NodeOffset int
	CommentMap ast.CommentMap
}

// Inspect run ast.Inspect for the file
func (f *FileInfo) Inspect(fun func(ast.Node) bool) {
	ast.Inspect(f.Ast, fun)
}

var newline = []byte{'\n'}

// GetNodeInfo returns *NodeInfo from ast.Node
func (f *FileInfo) GetNodeInfo(node ast.Node) *NodeInfo {
	lines := bytes.Split(f.Source, newline)
	start := int(node.Pos()-1) - f.NodeOffset
	end := int(node.End()) - f.NodeOffset
	sourceBeforeStart := f.Source[:start]
	numSourceBeforeStart := len(sourceBeforeStart)
	numLines := bytes.Count(sourceBeforeStart, newline) + 1
	return &NodeInfo{
		FilePath: f.Path,
		LineNum:  numLines,
		LineText: string(lines[numLines-1]),
		Pos:      start - numSourceBeforeStart,
		End:      end - numSourceBeforeStart,
	}
}

type NodeInfo struct {
	FilePath string
	LineNum  int
	LineText string
	Pos      int
	End      int
}

func (ni *NodeInfo) String() string {
	return fmt.Sprintf("%s [L%d](%d:%d) %s", ni.FilePath, ni.LineNum, ni.Pos, ni.End, ni.LineText)
}
