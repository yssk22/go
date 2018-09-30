package slice

import (
	"fmt"
	"reflect"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func assertSlice(v reflect.Value) {
	if k := v.Kind(); k != reflect.Slice {
		panic(fmt.Errorf("%s (%s) is not a slice", v.Type().Name(), k))
	}
}

func assertSliceFun(fType reflect.Type) {
	if fType.Kind() != reflect.Func {
		panic(fmt.Errorf("SliceFuncError: not a function"))
	}
	if fType.NumIn() != 2 {
		panic(fmt.Errorf("SliceFuncError: the second function must take two arguments"))
	}
	if fType.In(0).Kind() != reflect.Int {
		panic(fmt.Errorf("SliceFuncError: the second function must take int value on the first argument"))
	}
	if fType.In(1).Kind() == reflect.Struct {
		panic(fmt.Errorf(
			"SliceFuncError: the second function must not take struct value on the second argument, use %q instead",
			reflect.PtrTo(fType.In(1)),
		))
	}
}

const noErrorsInMultiError = "No errors"

// SliceError is an error collection as a single error.
// error[i] might be nil if there is no error.
type SliceError []error

// NewSliceError creates SliceError instance with the given size.
func NewSliceError(size int) SliceError {
	return SliceError(make([]error, size))
}

// Error implemnts error.Error()
func (se SliceError) Error() string {
	var firstError error
	var errorCount int
	for _, e := range se {
		if e != nil {
			if firstError == nil {
				firstError = e
			}
			errorCount++
		}
	}
	switch errorCount {
	case 0:
		return noErrorsInMultiError
	case 1:
		return firstError.Error()
	}
	return fmt.Sprintf("%s (and %d other errors)", firstError.Error(), errorCount)
}

// SplitByLength splits a list into multiple lists. Each lengths of lists must be up to `each`
// You can use the returned value as `[][]T` when you pass `[]T` for list.
func SplitByLength(list interface{}, eachSize int) interface{} {
	a := reflect.ValueOf(list)
	assertSlice(a)
	// create and allocate lists
	bucketType := a.Type()
	bucketListType := reflect.SliceOf(bucketType)
	tailSize := a.Len() % eachSize
	bucketListLen := a.Len() / eachSize
	if tailSize != 0 {
		bucketListLen++
	}
	bucketList := reflect.MakeSlice(bucketListType, bucketListLen, bucketListLen)

	// fill non-tail (must hold eachSize)
	for i := 0; i < bucketListLen-1; i++ {
		bucket := bucketList.Index(i)
		bucket.Set(reflect.MakeSlice(bucketType, eachSize, eachSize))
		offset := i * eachSize
		for j := 0; j < eachSize; j++ {
			bucket.Index(j).Set(a.Index(offset + j))
		}
	}

	// fill tail
	if tailSize == 0 {
		tailSize = eachSize
	}
	bucket := bucketList.Index(bucketListLen - 1)
	bucket.Set(reflect.MakeSlice(bucketType, tailSize, tailSize))
	offset := (bucketListLen - 1) * eachSize
	for j := 0; j < tailSize; j++ {
		bucket.Index(j).Set(a.Index(offset + j))
	}
	return bucketList.Interface()
}

// ToInterface converts []*T to []interface{}
func ToInterface(v interface{}) []interface{} {
	a := reflect.ValueOf(v)
	assertSlice(a)
	vv := make([]interface{}, a.Len())
	for i := range vv {
		vv[i] = a.Index(i).Interface()
	}
	return vv
}

// ToAddr converts []T to []*T
func ToAddr(v interface{}) interface{} {
	a := reflect.ValueOf(v)
	assertSlice(a)
	list := reflect.MakeSlice(
		reflect.SliceOf(reflect.PtrTo(a.Type().Elem())),
		a.Len(),
		a.Cap(),
	)
	for i := 0; i < a.Len(); i++ {
		list.Index(i).Set(a.Index(i).Addr())
	}
	return list.Interface()
}

// ToElem converts []*T to []T
func ToElem(v interface{}) interface{} {
	a := reflect.ValueOf(v)
	assertSlice(a)
	list := reflect.MakeSlice(
		reflect.SliceOf(a.Type().Elem().Elem()),
		a.Len(),
		a.Cap(),
	)
	for i := 0; i < a.Len(); i++ {
		list.Index(i).Set(a.Index(i).Elem())
	}
	return list.Interface()
}
