package iterator

import (
	"fmt"
	"reflect"
)

// import (
// 	"fmt"
// 	"reflect"
// 	"sync"
// )

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func assertSlice(v reflect.Value) {
	if k := v.Kind(); k != reflect.Slice {
		panic(fmt.Errorf("%s (%s) is not a slice", v.Type().Name(), k))
	}
}

func assertSliceFun(f reflect.Value) {
	fType := f.Type()
	if fType.NumIn() != 2 {
		panic(fmt.Errorf("ParallelError: the second function must take two arguments"))
	}
	if fType.In(0).Kind() != reflect.Int {
		panic(fmt.Errorf("ParallelError: the second function must take int value on the first argument"))
	}
	if fType.In(1).Kind() == reflect.Struct {
		panic(fmt.Errorf(
			"ParallelError: the second function must not take struct value on the second argument, use %q instead",
			reflect.PtrTo(fType.In(1)),
		))
	}
	if fType.NumOut() != 1 || !fType.Out(0).Implements(errorType) {
		panic(fmt.Errorf("ParallelError: the second function must return an error"))
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

// SplitSliceByLength splits a list into multiple lists. Each lengths of lists must be up to `each`
func SplitSliceByLength(list interface{}, eachSize int) interface{} {
	a := reflect.ValueOf(list)
	assertSlice(a)
	// create and allocate lists
	bucketType := a.Type()
	bucketListType := reflect.SliceOf(bucketType)
	tailSize := a.Len() % eachSize
	bucketListLen := a.Len()/eachSize + tailSize%2
	bucketList := reflect.MakeSlice(bucketListType, bucketListLen, bucketListLen)

	for i := 0; i < bucketListLen-1; i++ {
		bucket := bucketList.Index(i)
		bucket.Set(reflect.MakeSlice(bucketType, eachSize, eachSize))
		offset := i * eachSize
		for j := 0; j < eachSize; j++ {
			bucket.Index(j).Set(a.Index(offset + j))
		}
	}

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

// // ParallelSlice iterates a slice and execute `fun` in parallel
// func ParallelSlice(list interface{}, fun interface{}) error {
// 	a := reflect.ValueOf(list)
// 	assertSlice(a)
// 	f := reflect.ValueOf(fun)
// 	fType := f.Type()
// 	if fType.NumIn() != 2 {
// 		panic(fmt.Errorf("ParallelError: the second function must take two arguments"))
// 	}
// 	if fType.In(0).Kind() != reflect.Int {
// 		panic(fmt.Errorf("ParallelError: the second function must take int value on the first argument"))
// 	}
// 	if fType.In(1).Kind() == reflect.Struct {
// 		panic(fmt.Errorf(
// 			"ParallelError: the second function must not take struct value on the second argument, use %q instead",
// 			reflect.PtrTo(fType.In(1)),
// 		))
// 	}

// 	if fType.NumOut() != 1 || !fType.Out(0).Implements(errorType) {
// 		panic(fmt.Errorf("ParallelError: the second function must return an error"))
// 	}

// 	var wg sync.WaitGroup
// 	l := a.Len()
// 	shouldUsePtr := a.Type().Elem().Kind() == reflect.Struct
// 	errors := SliceError(make([]error, l))
// 	for i := 0; i < l; i++ {
// 		wg.Add(1)
// 		v := a.Index(i)
// 		if shouldUsePtr {
// 			v = v.Addr()
// 		}
// 		go func(i int, v reflect.Value) {
// 			defer wg.Done()
// 			defer func() {
// 				if x := recover(); x != nil {
// 					errors[i] = fmt.Errorf("%v", x)
// 				}
// 			}()
// 			out := f.Call([]reflect.Value{reflect.ValueOf(i), v})[0]
// 			if !out.IsNil() {
// 				errors[i] = out.Interface().(error)
// 			}
// 		}(i, v)
// 	}
// 	wg.Wait()
// 	any := false
// 	for i := 0; i < l; i++ {
// 		if errors[i] != nil {
// 			any = true
// 		}
// 	}
// 	if any {
// 		return errors
// 	}
// 	return nil
// }

// // MustParallelSlice is like ParallelSlice but panic if it returns an error
// func MustParallelSlice(list interface{}, fun interface{}) {
// 	if err := ParallelSlice(list, fun); err != nil {
// 		panic(err)
// 	}
// }

// // ParallelSliceWithMaxConcurrency is like ParallelSlice but spawn goroutines up to `n` concurrency
// func ParallelSliceWithMaxConcurrency(list interface{}, n int, fun interface{}) error {
// 	a1 := reflect.ValueOf(list)
// 	assertSlice(a1)
// 	f := reflect.ValueOf(fun)
// 	assertSliceFun(f)
// 	l := a1.Len()
// 	errors := SliceError(make([]error, l))

// 	if n > l {
// 		n = l
// 	}
// 	eachSize := l / n
// 	a := reflect.ValueOf(SplitByEach(list, eachSize))

// 	var wg sync.WaitGroup

// 	shouldUsePtr := a.Type().Elem().Kind() == reflect.Struct
// 	for i := 0; i < a.Len(); i++ {
// 		wg.Add(1)
// 		v := a.Index(i) // still slice as we used SplitByEach
// 		go func(i int, v reflect.Value) {
// 			defer wg.Done()
// 			for j := 0; j < v.Len(); j++ {
// 				idx := i*eachSize + j
// 				func() {
// 					defer func() {
// 						if x := recover(); x != nil {
// 							errors[idx] = fmt.Errorf("%v", x)
// 						}
// 					}()
// 					v1 := v.Index(j)
// 					if shouldUsePtr {
// 						v1 = v1.Addr()
// 					}
// 					out := f.Call([]reflect.Value{reflect.ValueOf(idx), v1})[0]
// 					if !out.IsNil() {
// 						errors[idx] = out.Interface().(error)
// 					}
// 				}()
// 			}
// 		}(i, v)
// 	}
// 	wg.Wait()
// 	any := false
// 	for i := 0; i < l; i++ {
// 		if errors[i] != nil {
// 			any = true
// 		}
// 	}
// 	if any {
// 		return errors
// 	}
// 	return nil
// }
