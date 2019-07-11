package slice

import (
	"fmt"
	"reflect"
)

// ToMap converts []T to map[key]T, where key is defined by fun
func ToMap(list interface{}, fun interface{}) interface{} {
	a := reflect.ValueOf(list)
	f := reflect.ValueOf(fun)
	fType := f.Type()
	assertSlice(a)
	assertSliceFun(fType)
	if fType.NumOut() != 1 {
		panic(fmt.Errorf("SliceFuncError: the second function must return an interface for the key"))
	}
	elemType := fType.In(1)
	keyType := fType.Out(0)
	m := reflect.MakeMap(reflect.MapOf(keyType, elemType))
	for i := 0; i < a.Len(); i++ {
		val := a.Index(i)
		key := f.Call([]reflect.Value{reflect.ValueOf(i), val})[0]
		m.SetMapIndex(key, val)
	}
	return m.Interface()
}
