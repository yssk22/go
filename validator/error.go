package validator

import (
	"bytes"
	"text/template"
)

// ValidationError is an error returned by Eval(Validator, interface{})
type ValidationError struct {
	Errors map[string][]*FieldError `json:"errors"`
	Value  interface{}              `json:"value"`
}

// FieldError represents an error in the validated object.
// It has a Message field which is a error message template,
// and Params map which would be applied to the error message.
type FieldError struct {
	Message string                 `json:"message"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// Returns string representation of the error
func (v *FieldError) String() string {
	var b bytes.Buffer
	t := template.Must(template.New("ve").Parse(v.Message))
	t.Execute(&b, v.Params)
	return b.String()
}

// Returns string representation of the error.
func (v *FieldError) Error() string {
	return v.String()
}

// NewFieldError creates a new ValidationError object.
func NewFieldError(message string, params map[string]interface{}) *FieldError {
	return &FieldError{
		Message: message,
		Params:  params,
	}
}
