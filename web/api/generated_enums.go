// Code generated by github.com/yssk22/go/generator DO NOT EDIT.
//
package api

import (
	"encoding/json"
	"fmt"
)

var (
	_RequestParameterFieldTypeValueToString = map[RequestParameterFieldType]string{
		RequestParameterFieldTypeBool:   "bool",
		RequestParameterFieldTypeInt:    "int",
		RequestParameterFieldTypeString: "string",
		RequestParameterFieldTypeFloat:  "float",
		RequestParameterFieldTypeTime:   "time",
		RequestParameterFieldTypeArray:  "array",
		RequestParameterFieldTypeObject: "object",
	}
	_RequestParameterFieldTypeStringToValue = map[string]RequestParameterFieldType{
		"bool":   RequestParameterFieldTypeBool,
		"int":    RequestParameterFieldTypeInt,
		"string": RequestParameterFieldTypeString,
		"float":  RequestParameterFieldTypeFloat,
		"time":   RequestParameterFieldTypeTime,
		"array":  RequestParameterFieldTypeArray,
		"object": RequestParameterFieldTypeObject,
	}
)

func (e RequestParameterFieldType) String() string {
	if str, ok := _RequestParameterFieldTypeValueToString[e]; ok {
		return str
	}
	return fmt.Sprintf("RequestParameterFieldType(%d)", e)
}

func (e RequestParameterFieldType) IsVaild() bool {
	_, ok := _RequestParameterFieldTypeValueToString[e]
	return ok
}

func ParseRequestParameterFieldType(s string) (RequestParameterFieldType, error) {
	if val, ok := _RequestParameterFieldTypeStringToValue[s]; ok {
		return val, nil
	}
	return RequestParameterFieldType(-1), fmt.Errorf("invalid value %q for RequestParameterFieldType", s)
}

func ParseRequestParameterFieldTypeOr(s string, or RequestParameterFieldType) RequestParameterFieldType {
	val, err := ParseRequestParameterFieldType(s)
	if err != nil {
		return or
	}
	return val
}

func MustParseRequestParameterFieldType(s string) RequestParameterFieldType {
	val, err := ParseRequestParameterFieldType(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (e RequestParameterFieldType) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _RequestParameterFieldTypeValueToString[e]; !ok {
		s = fmt.Sprintf("RequestParameterFieldType(%d)", e)
	}
	return json.Marshal(s)
}

func (e *RequestParameterFieldType) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := ParseRequestParameterFieldType(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*e = newval
	return nil
}

func (e *RequestParameterFieldType) Parse(s string) error {
	if val, ok := _RequestParameterFieldTypeStringToValue[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid value %q for RequestParameterFieldType", s)
}

var (
	_RequestParameterFormatValueToString = map[RequestParameterFormat]string{
		RequestParameterFormatQuery: "query",
		RequestParameterFormatForm:  "form",
		RequestParameterFormatJSON:  "json",
	}
	_RequestParameterFormatStringToValue = map[string]RequestParameterFormat{
		"query": RequestParameterFormatQuery,
		"form":  RequestParameterFormatForm,
		"json":  RequestParameterFormatJSON,
	}
)

func (e RequestParameterFormat) String() string {
	if str, ok := _RequestParameterFormatValueToString[e]; ok {
		return str
	}
	return fmt.Sprintf("RequestParameterFormat(%d)", e)
}

func (e RequestParameterFormat) IsVaild() bool {
	_, ok := _RequestParameterFormatValueToString[e]
	return ok
}

func ParseRequestParameterFormat(s string) (RequestParameterFormat, error) {
	if val, ok := _RequestParameterFormatStringToValue[s]; ok {
		return val, nil
	}
	return RequestParameterFormat(-1), fmt.Errorf("invalid value %q for RequestParameterFormat", s)
}

func ParseRequestParameterFormatOr(s string, or RequestParameterFormat) RequestParameterFormat {
	val, err := ParseRequestParameterFormat(s)
	if err != nil {
		return or
	}
	return val
}

func MustParseRequestParameterFormat(s string) RequestParameterFormat {
	val, err := ParseRequestParameterFormat(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (e RequestParameterFormat) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _RequestParameterFormatValueToString[e]; !ok {
		s = fmt.Sprintf("RequestParameterFormat(%d)", e)
	}
	return json.Marshal(s)
}

func (e *RequestParameterFormat) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := ParseRequestParameterFormat(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*e = newval
	return nil
}

func (e *RequestParameterFormat) Parse(s string) error {
	if val, ok := _RequestParameterFormatStringToValue[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid value %q for RequestParameterFormat", s)
}
