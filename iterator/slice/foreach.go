package slice

import (
	"fmt"
	"reflect"

	"github.com/yssk22/go/x/xerrors"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// ForEach iterate the slice
func ForEach(list interface{}, fun interface{}) error {
	a := reflect.ValueOf(list)
	f := reflect.ValueOf(fun)
	fType := f.Type()
	assertSlice(a)
	assertSliceFun(fType)
	if fType.NumOut() != 1 || fType.Out(0).Kind() != reflect.Interface || !fType.Out(0).Implements(errorInterface) {
		panic(fmt.Errorf("SliceFuncError: the second function must return an error"))
	}
	l := a.Len()
	err := xerrors.NewMultiError(l)
	for i := 0; i < a.Len(); i++ {
		out := f.Call([]reflect.Value{reflect.ValueOf(i), a.Index(i)})[0].Interface()
		if e, ok := out.(error); ok {
			err[i] = e
		}
	}
	if err.HasError() {
		return err
	}
	return nil
}
