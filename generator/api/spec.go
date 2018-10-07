package api

import (
	"fmt"
	"go/types"
	"strings"
)

// Spec represents API specification
type Spec struct {
	PathPattern string
	FuncName    string
	Method      string // one of Get, Post, Put, Delete
}

// Validate validates the spec
func (spec *Spec) Validate(pkg *types.Package) error {
	if obj := pkg.Scope().Lookup(spec.FuncName); obj == nil {
		return fmt.Errorf("%q is not found in %s", spec.FuncName, pkg.Name())
	}
	if len(spec.PathPattern) == 0 {
		return fmt.Errorf("path is empty")
	}
	switch spec.Method {
	case "Get", "Post", "Put", "Delete":
		break
	default:
		return fmt.Errorf("invalid")
	}
	return nil
}

func guessMethodByFunctionName(funcName string) string {
	if strings.HasPrefix(funcName, "get") {
		return "Get"
	}
	if strings.HasPrefix(funcName, "update") {
		return "Put"
	}
	if strings.HasPrefix(funcName, "create") {
		return "Post"
	}
	if strings.HasPrefix(funcName, "delete") {
		return "Delete"
	}
	return "" // unknown
}
