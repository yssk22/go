package flow

import (
	"fmt"
	"strings"
)

// Spec represents type specification
type Spec struct {
	TypeName string
	FlowType FlowType
}

type Field struct {
	Name string
	Type FlowType
}

type FlowType interface {
	GetExpr() string
}

type FlowTypePrimitive string

func (f FlowTypePrimitive) GetExpr() string {
	return string(f)
}

// Primitive flow types
const (
	FlowTypeString FlowTypePrimitive = "string"
	FlowTypeNumber FlowTypePrimitive = "number"
	FlowTypeBool   FlowTypePrimitive = "boolean"
	FlowTypeDate   FlowTypePrimitive = "Date"
	FlowTypeAny    FlowTypePrimitive = "any"
)

type FlowTypeMaybe struct {
	Elem FlowType
}

func (f *FlowTypeMaybe) GetExpr() string {
	return fmt.Sprintf("?%s", f.Elem.GetExpr())
}

type FlowTypeArray struct {
	ElemenetType FlowType
}

func (f *FlowTypeArray) GetExpr() string {
	return fmt.Sprintf("Array<%s>", f.ElemenetType.GetExpr())
}

type FlowTypeObject struct {
	Fields []FlowTypeObjectField
}

func (f *FlowTypeObject) GetExpr() string {
	var lines []string
	lines = append(lines, "{")
	for _, field := range f.Fields {
		if field.OmitEmpty {
			lines = append(
				lines,
				fmt.Sprintf("%s?: %s,", field.Name, field.Type.GetExpr()),
			)
		} else {
			lines = append(
				lines,
				fmt.Sprintf("%s: %s,", field.Name, field.Type.GetExpr()),
			)
		}
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

type FlowTypeObjectField struct {
	Name      string
	Type      FlowType
	OmitEmpty bool
}

// Other Named object
type FlowTypeNamed struct {
	Name       string
	ImportPath string
	ImportName string
}

func (f *FlowTypeNamed) GetExpr() string {
	if f.ImportName != "" {
		return fmt.Sprintf("%s.%s", f.ImportName, f.Name)
	}
	return f.Name
}
