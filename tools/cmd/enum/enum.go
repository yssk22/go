package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	exact "go/constant"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/speedland/go/x/xstrings"
)

const defaultOutput = "enum_helper.go"

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_enum.go.go")
)

func main() {
	log.SetPrefix("[enum] ")
	log.SetFlags(0)
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, directory := range args {
		var g Generator
		g.parsePackageDir(directory)
		g.Printf("// Code generated by \"enum %s\"; DO NOT EDIT\n", strings.Join(os.Args[1:], " "))
		g.Printf("\n")
		g.Printf("package %s", g.pkg.name)
		g.Printf("\n")
		g.Printf("import (\n")
		g.Printf("\t\"encoding/json\"\n")
		g.Printf("\t\"fmt\"\n")
		g.Printf(")\n")
		genCount := 0
		for _, typeName := range types {
			if g.generate(typeName) {
				genCount++
			}
		}
		if genCount == 0 {
			log.Printf("No enum found.")
			continue
		}
		src := g.format()
		outputName := *output
		if outputName == "" {
			outputName = filepath.Join(directory, fmt.Sprintf("%s_enum.go", xstrings.ToSnakeCase(types[0])))
		}
		err := ioutil.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}
}

type Generator struct {
	buff bytes.Buffer
	pkg  *Package
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buff, format, args...)
}

// parsePackageDir parses the package residing in the directory.
func (g *Generator) parsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}
	var allFiles = [][]string{pkg.GoFiles, pkg.CgoFiles, pkg.SFiles}
	var names []string
	for _, files := range allFiles {
		for _, file := range files {
			names = append(names, filepath.Join(directory, file))
		}
	}
	g.parsePackage(directory, names, nil)
}

func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, text, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	// Type check the package.
	g.pkg.check(fs, astFiles)
}

func (g *Generator) generate(typeName string) bool {
	var values []Value
	for _, file := range g.pkg.files {
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}
	if len(values) == 0 {
		return false
	}
	runs := groupByTypes(values)
	g.buildRuns(runs)
	return true
}

func (g *Generator) format() []byte {
	src, err := format.Source(g.buff.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buff.Bytes()
	}
	return src
}

// buildMultipleRuns generates the variables and String method for multiple runs of contiguous values.
// For this pattern, a single Printf format won't do.
func (g *Generator) buildRuns(runs [][]Value) {
	g.Printf("\n")
	g.declareMaps(runs)
	g.Printf("\n")
	for _, values := range runs {
		log.Printf("\t%s (%d values)", values[0].typ, len(values))
		for _, funcTemplate := range supportedFuncs {
			g.Printf(funcTemplate, values[0].typ)
			g.Printf("\n")
		}
	}
}

func (g *Generator) declareMaps(runs [][]Value) {
	g.Printf("var (\n")
	for _, values := range runs {
		typ := values[0].typ
		g.Printf("\t_%sValueToString = map[%s]string{\n", typ, typ)
		for _, v := range values {
			g.Printf("\t\t%s: %q,\n", v.name, v.String())
		}
		g.Printf("\t}\n")
		g.Printf("\t_%sStringToValue = map[string]%s{\n", typ, typ)
		for _, v := range values {
			g.Printf("\t\t%q: %v,\n", v.String(), v.name)
		}
		g.Printf("\t}\n")
		g.Printf("\n")
	}
	g.Printf(")\n")
}

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	files    []*File
	typesPkg *types.Package
}

// check type-checks the package. The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	values   []Value // Accumulator for constant values of that type.
	typeName string
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value, the constant is untyped.
			// Skip this vspec and reset the remembered type.
			typ = ""
			continue
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeName {
			// This is not the type we're looking for.
			continue
		}
		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
			}
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != exact.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			i64, isInt := exact.Int64Val(value)
			u64, isUint := exact.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
			}
			if !isInt {
				u64 = uint64(i64)
			}
			str, err := nameToString(name.Name, typ)
			if err != nil {
				log.Fatalf("%v", err)
			}
			v := Value{
				name:     name.Name,
				typ:      typ,
				value:    u64,
				signed:   info&types.IsUnsigned == 0,
				strValue: str,
			}
			f.values = append(f.values, v)
		}
	}
	return false
}

func nameToString(name string, typeName string) (string, error) {
	if !strings.HasPrefix(name, typeName) {
		return "", fmt.Errorf("Naming violation: %s must starts with %s", name, typeName)
	}
	return xstrings.ToSnakeCase(strings.TrimPrefix(name, typeName)), nil
}

// Value represents a declared constant.
type Value struct {
	name     string // The name of the constant.
	typ      string // Type type of the constant
	value    uint64 // Will be converted to int64 when needed.
	signed   bool   // Whether the constant is a signed type.
	strValue string // The string representation for the name.
}

func (v *Value) String() string {
	return v.strValue
}

// byValue lets us sort the constants into increasing order.
// We take care in the Less method to sort in signed or unsigned order,
// as appropriate.
type byValue []Value

func (b byValue) Len() int      { return len(b) }
func (b byValue) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byValue) Less(i, j int) bool {
	if b[i].signed {
		return int64(b[i].value) < int64(b[j].value)
	}
	return b[i].value < b[j].value
}

func groupByTypes(values []Value) [][]Value {
	// We use stable sort so the lexically first name is chosen for equal elements.
	sort.Stable(byValue(values))

	var typIndex = make(map[string]int)
	var runs [][]Value
	idx := 0
	for _, v := range values {
		if _, ok := typIndex[v.typ]; !ok {
			runs = append(runs, []Value{})
			typIndex[v.typ] = idx
			idx++
		}
		runs[typIndex[v.typ]] = append(runs[typIndex[v.typ]], v)
	}
	return runs
}

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}
