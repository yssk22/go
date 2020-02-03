package slice

import (
	"fmt"
	"reflect"

	"github.com/yssk22/go/x/xerrors"
)

// Map iterate the slice
// list must be []T, where fun must be func(i int, v T) (T1, error)
func Map(list interface{}, fun interface{}) (interface{}, error) {
	a := reflect.ValueOf(list)
	f := reflect.ValueOf(fun)
	fType := f.Type()
	assertSlice(a)
	assertSliceFun(fType)
	if fType.NumOut() != 2 || fType.Out(1).Kind() != reflect.Interface || !fType.Out(1).Implements(errorInterface) {
		panic(fmt.Errorf("SliceFuncError: the second function must return (interface{}, error)"))
	}
	l := a.Len()

	shouldUsePtr := a.Type().Elem().Kind() == reflect.Struct
	mapped := reflect.MakeSlice(reflect.SliceOf(fType.Out(0)), l, l)
	err := xerrors.NewMultiError(l)
	for i := 0; i < a.Len(); i++ {
		v := a.Index(i)
		if shouldUsePtr {
			v = v.Addr()
		}
		out := f.Call([]reflect.Value{reflect.ValueOf(i), v})
		mapped.Index(i).Set(out[0])
		e := out[1].Interface()
		if _, ok := e.(error); ok {
			err[i] = e.(error)
		}
	}
	return mapped.Interface(), err.ToReturn()
}
