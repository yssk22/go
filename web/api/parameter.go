package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/yssk22/go/x/xtime"

	"github.com/yssk22/go/web/response"
)

// RequestParameterFormat is an enum to represent the format of request parameters
// @enum
type RequestParameterFormat int

// Available RequestParameterFormat values
const (
	RequestParameterFormatQuery RequestParameterFormat = iota
	RequestParameterFormatForm
	RequestParameterFormatJSON
)

// RequestParameterFieldSpec is a spec to parse the parameters
type RequestParameterFieldSpec struct {
	Type     RequestParameterFieldType
	Default  interface{}
	Required bool
}

// RequestParameterFieldType is a type enum for request parameters
// @enum
type RequestParameterFieldType int

// Available RequestParameterFieldType values
const (
	RequestParameterFieldTypeBool RequestParameterFieldType = iota
	RequestParameterFieldTypeInt
	RequestParameterFieldTypeString
	RequestParameterFieldTypeFloat
	RequestParameterFieldTypeTime
	RequestParameterFieldTypeArray
	RequestParameterFieldTypeObject
)

// ParameterParser is to parse parameter value
type ParameterParser struct {
	specs  map[string]*RequestParameterFieldSpec
	format RequestParameterFormat
}

// NewParameterParser return a new NewParameterParser instance
func NewParameterParser(format RequestParameterFormat) *ParameterParser {
	return &ParameterParser{
		specs:  make(map[string]*RequestParameterFieldSpec),
		format: format,
	}
}

// Type to set the parameter field type given by the key
func (pp *ParameterParser) Type(key string, t RequestParameterFieldType) *ParameterParser {
	pp.spec(key).Type = t
	return pp
}

// Default to set the default value for the parameter field given by the key
func (pp *ParameterParser) Default(key string, val interface{}) *ParameterParser {
	pp.spec(key).Default = val
	return pp
}

// Required to set the parameter field given by the key required
func (pp *ParameterParser) Required(key string) *ParameterParser {
	pp.spec(key).Required = true
	return pp
}

func (pp *ParameterParser) spec(key string) *RequestParameterFieldSpec {
	spec, ok := pp.specs[key]
	if !ok {
		spec = &RequestParameterFieldSpec{}
		pp.specs[key] = spec
	}
	return spec
}

// Parse runs to parse paraemters in the request
func (pp *ParameterParser) Parse(req *http.Request, v interface{}) *Error {
	var err error
	var jsonBytes []byte // normalized JSON bytes
	switch pp.format {
	case RequestParameterFormatQuery:
		jsonBytes, err = pp.toJSONString(req.URL.Query())
		break
	case RequestParameterFormatForm:
		req.ParseForm()
		jsonBytes, err = pp.toJSONString(req.PostForm)
		break
	case RequestParameterFormatJSON:
		jsonBytes, err = ioutil.ReadAll(req.Body)
	}
	if err != nil {
		if _, ok := err.(*Error); ok {
			return err.(*Error)
		}
		return ServerError
	}
	err = json.Unmarshal(jsonBytes, v)
	if err != nil {
		return ServerError
	}

	return nil
}

func (pp *ParameterParser) toJSONString(v url.Values) ([]byte, error) {
	var fe = newFieldErrors()
	var err error
	var m = make(map[string]interface{})
	for k, spec := range pp.specs {
		val, ok := v[k]
		hasKey := ok && len(val) > 0
		if !hasKey {
			if spec.Required {
				fe.addError(k, fmt.Errorf("required"))
			}
			if spec.Default != nil {
				m[k] = spec.Default
			}
			continue
		}
		strValue := val[0]
		switch spec.Type {
		case RequestParameterFieldTypeString:
			m[k] = strValue
		case RequestParameterFieldTypeBool:
			m[k] = strValue == "1" || strValue == "true"
		case RequestParameterFieldTypeInt:
			if m[k], err = strconv.Atoi(strValue); err != nil {
				fe.addError(k, fmt.Errorf("must be int, but %q", strValue))
			}
			break
		case RequestParameterFieldTypeFloat:
			if m[k], err = strconv.ParseFloat(strValue, 64); err != nil {
				fe.addError(k, fmt.Errorf("must be float, but %q", strValue))
			}
			break
		case RequestParameterFieldTypeTime:
			if m[k], err = xtime.Parse(strValue); err != nil {
				if m[k], err = xtime.ParseDateDefault(strValue); err != nil {
					fe.addError(k, fmt.Errorf("must be time format, but %q", strValue))
				}
			}
			break
		case RequestParameterFieldTypeArray:
			var mm []interface{}
			if err = json.Unmarshal([]byte(strValue), &mm); err != nil {
				fe.addError(k, fmt.Errorf("must be json array, but %q", strValue))
			} else {
				m[k] = mm
			}
			break
		case RequestParameterFieldTypeObject:
			log.Println("object")
			log.Println(k, strValue)
			var mm = make(map[string]interface{})
			if err = json.Unmarshal([]byte(strValue), &mm); err != nil {
				fe.addError(k, fmt.Errorf("must be json object, but %q", strValue))
			} else {
				m[k] = mm
			}
			break
		}
	}
	if err = fe.ToError(); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type fieldErrors struct {
	errors map[string][]string
}

func newFieldErrors() *fieldErrors {
	return &fieldErrors{
		errors: make(map[string][]string),
	}
}

func (fe *fieldErrors) addError(key string, e error) {
	if _, ok := fe.errors[key]; !ok {
		fe.errors[key] = make([]string, 0)
	}
	fe.errors[key] = append(fe.errors[key], e.Error())
}

func (fe *fieldErrors) ToError() error {
	if len(fe.errors) == 0 {
		return nil
	}
	return (&Error{
		Code:    "invalid_parameter",
		Message: "one or more parameters are invalid",
		Extra: map[string]interface{}{
			"errors": fe.errors,
		},
		Status: response.HTTPStatusBadRequest,
	})
}
