package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yssk22/go/x/xerrors"

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
	Type     RequestParameterFieldType `json:"type"`
	Default  interface{}               `json:"default"`
	Required bool                      `json:"required"` // this checks the parameter has a key or not.
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

// ValueOf returns a typed value for strValue
func (t RequestParameterFieldType) ValueOf(strValue string) (interface{}, error) {
	switch t {
	case RequestParameterFieldTypeString:
		return strValue, nil
	case RequestParameterFieldTypeBool:
		return (strValue == "1" || strValue == "true"), nil
	case RequestParameterFieldTypeInt:
		if v, err := strconv.Atoi(strValue); err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("must be int, but %q", strValue)
	case RequestParameterFieldTypeFloat:
		if v, err := strconv.ParseFloat(strValue, 64); err == nil {
			return v, nil
		}
		return 0.0, fmt.Errorf("must be float, but %q", strValue)
	case RequestParameterFieldTypeTime:
		if v, err := xtime.Parse(strValue); err == nil {
			return v, nil
		}
		if v, err := xtime.ParseDateDefault(strValue); err == nil {
			return v, nil
		}
		return time.Time{}, fmt.Errorf("must be time format, but %q", strValue)
	case RequestParameterFieldTypeArray:
		var mm []interface{}
		if err := json.Unmarshal([]byte(strValue), &mm); err == nil {
			return mm, nil
		}
		return nil, fmt.Errorf("must be json array, but %q", strValue)
	case RequestParameterFieldTypeObject:
		var mm = make(map[string]interface{})
		if err := json.Unmarshal([]byte(strValue), &mm); err == nil {
			return mm, nil
		}
		return nil, fmt.Errorf("must be json object, but %q", strValue)
	}
	return nil, fmt.Errorf("unknown type %q to evaluate %q", t, strValue)
}

// ParameterParser is to parse parameter value
type ParameterParser struct {
	Specs  map[string]*RequestParameterFieldSpec `json:"specs"`
	Format RequestParameterFormat                `json:"format"`
}

// NewParameterParser return a new NewParameterParser instance
func NewParameterParser(format RequestParameterFormat) *ParameterParser {
	return &ParameterParser{
		Specs:  make(map[string]*RequestParameterFieldSpec),
		Format: format,
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
	spec, ok := pp.Specs[key]
	if !ok {
		spec = &RequestParameterFieldSpec{}
		pp.Specs[key] = spec
	}
	return spec
}

// Parse runs to parse paraemters in the request
func (pp *ParameterParser) Parse(req *http.Request, v interface{}) *Error {
	var err error
	var jsonBytes []byte // normalized JSON bytes
	switch pp.Format {
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
		xerrors.MustNil(err)
	}
	err = json.Unmarshal(jsonBytes, v)
	xerrors.MustNil(err)
	if validatable, ok := v.(Validatable); ok {
		errors := NewFieldErrorCollection()
		if err = validatable.Validate(req.Context(), errors); err != nil {
			return BadRequest
		}
		if err = errors.ToError(); err != nil {
			return err.(*Error)
		}
	}
	return nil
}

func (pp *ParameterParser) toJSONString(v url.Values) ([]byte, error) {
	var err error
	var errors = NewFieldErrorCollection()
	var m = make(map[string]interface{})
	for k, spec := range pp.Specs {
		val, ok := v[k]
		hasKey := ok && len(val) > 0
		if !hasKey {
			if spec.Required {
				errors.Add(k, fmt.Errorf("required"))
			}
			if spec.Default != nil {
				m[k] = spec.Default
			}
			continue
		}
		if v, err := spec.Type.ValueOf(val[0]); err != nil {
			errors.Add(k, err)
		} else {
			m[k] = v
		}
	}
	if err = errors.ToError(); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// FieldErrorCollection is a struct to collect field errors
type FieldErrorCollection map[string][]string

// NewFieldErrorCollection returns a new FieldErrorCollection
func NewFieldErrorCollection() FieldErrorCollection {
	return FieldErrorCollection(make(map[string][]string))
}

// Add to add an error associated the given field key
func (collection FieldErrorCollection) Add(key string, errors ...error) {
	var errList []string
	for _, e := range errors {
		if e != nil {
			errList = append(errList, e.Error())
		}
	}
	if len(errList) == 0 {
		return
	}
	if _, ok := collection[key]; !ok {
		collection[key] = make([]string, 0)
	}
	collection[key] = append(collection[key], errList...)
}

// ToError converts the collection to *Error object
func (collection FieldErrorCollection) ToError() error {
	if len(collection) == 0 {
		return nil
	}
	return (&Error{
		Code:    "invalid_parameter",
		Message: "one or more parameters are invalid",
		Extra: map[string]interface{}{
			"errors": collection,
		},
		Status: response.HTTPStatusBadRequest,
	})
}
