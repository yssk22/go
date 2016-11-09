package validator

import (
	"reflect"
	"regexp"

	"github.com/speedland/go/x/xreflect"
)

var requiredError = NewFieldError("must be required.", nil)
var requiredFunc = func(v interface{}) *FieldError {
	if s, ok := v.(string); ok {
		if s == "" {
			return requiredError
		}
	} else if s, ok := v.([]byte); ok {
		if len(s) == 0 {
			return requiredError
		}
	} else if xreflect.IsNil(v) {
		return requiredError
	}
	return nil
}

var minFunc = func(n int64) Func {
	var minError = NewFieldError(
		"must be more than or equal to {{.min}}",
		map[string]interface{}{
			"min": n,
		},
	)
	return func(v interface{}) *FieldError {
		if asInt(v) >= n {
			return nil
		}
		return minError
	}
}

var maxFunc = func(n int64) Func {
	var maxError = NewFieldError(
		"must be less than or equal to {{.max}}",
		map[string]interface{}{
			"max": n,
		},
	)
	return func(v interface{}) *FieldError {
		if asInt(v) <= n {
			return nil
		}
		return maxError
	}
}

var matchFunc = func(str string) Func {
	var matchError = NewFieldError(
		"not match with '{{.regexp}}'",
		map[string]interface{}{
			"regexp": str,
		},
	)
	var exp = regexp.MustCompile(str)
	return func(v interface{}) *FieldError {
		if isMatched(v, exp) {
			return nil
		}
		return matchError
	}
}

var unmatchFunc = func(str string) Func {
	var unmatchError = NewFieldError(
		"match with '{{.regexp}}'",
		map[string]interface{}{
			"regexp": str,
		},
	)
	var exp = regexp.MustCompile(str)
	return func(v interface{}) *FieldError {
		if !isMatched(v, exp) {
			return nil
		}
		return unmatchError
	}
}

func asInt(v interface{}) int64 {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return int64(val.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	case reflect.Map, reflect.Slice, reflect.Array:
		if val.IsNil() {
			return 0
		}
		return int64(val.Len())
	}
	return 0
}

func isMatched(v interface{}, with *regexp.Regexp) bool {
	switch v.(type) {
	case string:
		if with.MatchString(v.(string)) {
			return true
		}
	case []byte:
		if with.Match(v.([]byte)) {
			return true
		}
	}
	return false
}
