package slice

import (
	"fmt"
	"reflect"
)

// Comparable is an interface to compare a slice element with x
type Comparable interface {
	Compare(interface{}) (int, error)
}

// ErrInvalidComparison is an error object returend by Compare if not comparable.
var ErrInvalidComparison = fmt.Errorf("invalid comparison")

// MustCompare returns the compared result with v1 and v2. 0: v1 == v2, 1: v1 > v2, -1: v1 < v2.
func MustCompare(v1 Comparable, v2 interface{}) int {
	d, err := v1.Compare(v2)
	if err != nil {
		panic(err)
	}
	return d
}

// AppendIfMissing returns a new list of `list` if v doesn't present.
func AppendIfMissing(list interface{}, v ...interface{}) interface{} {
	a := reflect.ValueOf(list)
	assertSlice(a)
	l := a.Len()
	for _, vv := range v {
		var comparable, isComparable = vv.(Comparable)
		var exists = false
		for i := 0; i < l; i++ {
			val := a.Index(i).Interface()
			if val == vv {
				exists = true
				break
			}
			if isComparable {
				if MustCompare(comparable, val) == 0 {
					exists = true
					break
				}
			}
		}
		if !exists {
			a = reflect.Append(a, reflect.ValueOf(vv))
		}
	}
	return a.Interface()
}
