package slice

import (
	"fmt"
	"reflect"
)

// Filter returns a new filtered list
// when the list is []T, fun must be func(i int, v T) bool, then it returns elements where func(i, v) returns true.
func Filter(list interface{}, fun interface{}) interface{} {
	a := reflect.ValueOf(list)
	f := reflect.ValueOf(fun)
	fType := f.Type()
	assertSlice(a)
	assertSliceFun(fType)
	if fType.NumOut() != 1 || fType.Out(0).Kind() != reflect.Bool {
		panic(fmt.Errorf("SliceFuncError: the second function must return an boolr"))
	}
	l := a.Len()
	var filteredValues []reflect.Value
	for i := 0; i < a.Len(); i++ {
		v := a.Index(i)
		out := f.Call([]reflect.Value{reflect.ValueOf(i), v})[0].Bool()
		if !out {
			filteredValues = append(filteredValues, v)
		}
	}
	l = len(filteredValues)
	filtered := reflect.MakeSlice(a.Type(), l, l)
	for i, v := range filteredValues {
		filtered.Index(i).Set(v)
	}
	return filtered.Interface()
}
