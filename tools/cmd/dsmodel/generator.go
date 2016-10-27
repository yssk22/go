package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"html/template"
	"io"
)

// Generator implements generator.Generator
type Generator struct {
	Package string
	Type    string
	Fields  *ast.FieldList
}

// Inspect implements Generator#Inspect
func (g *Generator) Inspect(node ast.Node) bool {
	if g.Fields != nil && g.Package != "" {
		// no more inspection needed.
		return true
	}
	switch node.(type) {
	case *ast.File:
		g.Package = node.(*ast.File).Name.Name
		return true
	case *ast.GenDecl:
		decl := node.(*ast.GenDecl)
		return g.inspecgGenDecl(decl)
	default:
		return true
	}
}

// GenSource implements Generator#GenSource
func (g *Generator) GenSource(w io.Writer) error {
	if g.Fields == nil {
		return fmt.Errorf("no struct to be generated")
	}
	t := template.Must(template.New("template").Parse(codeTemplate))
	return t.Execute(w, g)
}

func (g *Generator) inspecgGenDecl(decl *ast.GenDecl) bool {
	if decl.Tok != token.TYPE {
		return true
	}
	for _, spec := range decl.Specs {
		t, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		structType, ok := t.Type.(*ast.StructType)
		if !ok {
			continue
		}
		if t.Name.Name != g.Type {
			continue
		}
		g.Fields = structType.Fields
		return true
	}
	return true
}
