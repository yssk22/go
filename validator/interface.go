package validator

import (
	"fmt"
	"reflect"
)

// Required validate v is not nil
func Required(v interface{}) error {
	if v == nil {
		return fmt.Errorf("required")
	}
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr && vv.IsNil() {
		return fmt.Errorf("required")
	}
	return nil
}
