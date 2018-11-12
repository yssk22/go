package api

import (
	"fmt"

	"github.com/yssk22/go/web/api"

	"github.com/yssk22/go/generator"
)

type requestMethod string

const (
	requestMethodGet     = "Get"
	requestMethodPost    = "Post"
	requestMethodPut     = "Put"
	requestMethodDelete  = "Delete"
	requestMethodUnknown = "Unknown"
)

// Spec represents API specification
type Spec struct {
	PathPattern         string
	PathParameters      []string
	StructuredParameter *StructuredParameter
	FuncName            string
	Method              requestMethod
}

// ParameterType represents the type information for a parameter
type ParameterType struct {
	Name         string
	Package      string
	PackageAlias string
}

func (pt *ParameterType) String() string {
	s := pt.PackageAlias
	if s == "" {
		s = pt.Package
	}
	if s != "" {
		s = s + "."
	}
	return fmt.Sprintf("%s%s", s, pt.Name)
}

// ResolveAlias resolves PackageAlias field with the given Dependency object.
func (pt *ParameterType) ResolveAlias(d *generator.Dependency) {
	if pt.Package != "" {
		pt.PackageAlias = d.Add(pt.Package)
	}
}

type StructuredParameter struct {
	Type   *ParameterType
	Parser *api.ParameterParser
}
