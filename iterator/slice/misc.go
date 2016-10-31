package slice

import "reflect"

// ToInterfaceSlice converts []T to []interface{}
// See https://github.com/golang/go/wiki/InterfaceSlice
func ToInterfaceSlice(list interface{}) []interface{} {
	a := reflect.ValueOf(list)
	assertSlice(a)
	size := a.Len()
	ia := make([]interface{}, size, size)
	for i := range ia {
		ia[i] = a.Index(i).Interface()
	}
	return ia
}
