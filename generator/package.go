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
	"regexp"
	"strings"

	"github.com/yssk22/go/x/xerrors"
)

// Signature represents @signature key=value key=value ... and corresponding ast.Node
type Signature struct {
	Name   string
	Params map[string]string
	Node   ast.Node
	Source *FileInfo
}

// GenError returns an wrapped error object with the node information
// n must be the ancestor of Signature node or nil
func (s *Signature) GenError(e error, n ast.Node) error {
	i := s.Source.GetNodeInfo(s.Node)
	if n != nil {
		ii := s.Source.GetNodeInfo(n)
		return xerrors.Wrap(e, ii.String())
	}
	return xerrors.Wrap(e, i.String())
}

// PackageInfo is a package infomaiton that generator is analysing
type PackageInfo struct {
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
			Ast:        parsedFile,
			CommentMap: ast.NewCommentMap(fs, parsedFile, parsedFile.Comments),
		})
		parsedFiles = append(parsedFiles, parsedFile)
	}
	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{},
		Defs:  map[*ast.Ident]types.Object{},
		Uses:  map[*ast.Ident]types.Object{},
	}
	pkg, err := conf.Check(".", fs, parsedFiles, info)
	if err != nil {
		return nil, xerrors.Wrap(err, "type check error -- you may need to run `go install ./...` at first")
	}
	pkg.SetName(importedPackage.Name)
	return &PackageInfo{
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

// CollectSignatures returns all `@s key=value ..` signatures
func (p *PackageInfo) CollectSignatures(s string) []*Signature {
	var signatures []*Signature
	var re = regexp.MustCompile(fmt.Sprintf("^\\s*@%s\\s*", s))
	for _, f := range p.Files {
		for node, commentGroups := range f.CommentMap {
			for _, c := range commentGroups {
				for _, line := range c.List {
					idx := re.FindStringIndex(line.Text)
					if len(idx) > 0 {
						remains := line.Text[idx[0]+len(s)+1:]
						params := parseSignatureParams(remains)
						signatures = append(signatures, &Signature{
							Name:   s,
							Params: params,
							Node:   node,
							Source: f,
						})
					}
				}
			}
		}
	}
	return signatures
}

func parseSignatureParams(s string) map[string]string {
	params := make(map[string]string)
	arguments := strings.Split(s, " ")
	for _, arg := range arguments {
		var key, value string
		idx := strings.Index(arg, "=")
		if idx < 0 {
			key = s
			value = ""
		} else {
			key = arg[:idx]
			value = arg[idx+1:]
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" {
			params[key] = value
		}
	}
	return params
}

// FileInfo is a file info that generator is analysing
type FileInfo struct {
	Path       string
	Ast        *ast.File
	Source     []byte
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
	start := int(node.Pos() - 1)
	end := int(node.End())
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
