package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/types"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yssk22/go/x/xerrors"
	"golang.org/x/tools/go/packages"
)

// PackageInfo is a package infomaiton that generator is analysing
type PackageInfo struct {
	Name     string
	Package  *types.Package
	TypeInfo *types.Info
	Files    []*FileInfo
}

// resolveGoImportPath resolves directry path to go import statement path.
func resolveGoImportPath(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", xerrors.Wrap(err, "filepath.Abs(%q) returns an error", build.Default.GOPATH)
	}
	absModuleRootPath, moduleName := findGoModInfo(absPath)
	if absModuleRootPath == "" {
		// go.mod not found
		absGoPath, err := filepath.Abs(build.Default.GOPATH)
		if err != nil {
			return "", xerrors.Wrap(err, "filepath.Abs(%q) returns an error", build.Default.GOPATH)
		}
		absGoPath = filepath.Join(absGoPath, "src")
		if !strings.HasPrefix(absPath, absGoPath) {
			return "", fmt.Errorf("not in $GOPATH/src (%s)", absGoPath)
		}
		offset := len(absGoPath) + 1
		return absPath[offset:], nil
	}

	// absPath is ${absModuleRootPath}/my/package/path/dir
	// and the import path should be ${moduleName}/my/package/path/dir
	return path.Join(moduleName, strings.TrimPrefix(absPath, absModuleRootPath)), nil
}

var (
	moduleDefRe = regexp.MustCompile("module\\s+(\\S+)\n")
)

// findGoModInfo finds go.mod file under the directory (or parent directories)
// and resturns the directory path where go.mod exists and module string declared in go.mod file.
func findGoModInfo(dir string) (string, string) {
	gomod := filepath.Join(dir, "go.mod")
	_, err := os.Stat(gomod)
	if err != nil {
		if os.IsNotExist(err) {
			if dir == "/" {
				return "", ""
			}
			return findGoModInfo(filepath.Join(dir, "..") + "/")
		}
		panic(err)
	}
	contents, err := ioutil.ReadFile(gomod)
	xerrors.MustNil(err)
	found := moduleDefRe.Copy().FindSubmatch(contents)
	if len(found) == 0 {
		panic(fmt.Errorf("could not find module declaration in %s", gomod))
	}
	return dir, string(found[1])
}

func parsePackage(dir string) (*PackageInfo, error) {
	resolvedGoImportPath, err := resolveGoImportPath(dir)
	if err != nil {
		return nil, err
	}
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  dir,
	}
	pkgs, err := packages.Load(cfg, resolvedGoImportPath)
	if err != nil {
		return nil, xerrors.Wrap(err, "could not load the package in %s", dir)
	}
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, xerrors.Wrap(pkg.Errors[0], "could not load the package in %s", dir)
	}
	if len(pkg.GoFiles) == 0 {
		return nil, fmt.Errorf("no go file is parsed in %s", dir)
	}
	// fs := token.NewFileSet()
	var files []*FileInfo
	for i, path := range pkg.GoFiles {
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, xerrors.Wrap(err, "io error: %s", path)
		}

		parsedFile := pkg.Syntax[i]
		// log.Println("File", path, parsedFile.Pos(), parsedFile.End(), offset, parsedFile.End()-parsedFile.Pos(), len(src))
		files = append(files, &FileInfo{
			Path:       path,
			Source:     src,
			Ast:        parsedFile,
			CommentMap: ast.NewCommentMap(pkg.Fset, parsedFile, parsedFile.Comments),
		})
	}
	// NOTE: We use *PackageInfo when we used in go 1.8.
	// We may want to return pkg directly in the future.
	return &PackageInfo{
		Name:     pkg.Name,
		Package:  pkg.Types,
		TypeInfo: pkg.TypesInfo,
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
	CommentMap ast.CommentMap
}

// Inspect run ast.Inspect for the file
func (f *FileInfo) Inspect(fun func(ast.Node) bool) {
	ast.Inspect(f.Ast, fun)
}

var newline = []byte{'\n'}

// GetNodeInfo returns *NodeInfo from ast.Node
func (f *FileInfo) GetNodeInfo(node ast.Node) *NodeInfo {
	fileOffset := int(f.Ast.Pos())
	lines := bytes.Split(f.Source, newline)
	start := int(node.Pos()) - fileOffset
	end := int(node.End()) - fileOffset
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
