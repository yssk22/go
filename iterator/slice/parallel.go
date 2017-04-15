package slice

import (
	"fmt"
	"reflect"
	"sync"
)

// ParallelOption is an option for Parallel* methods
type ParallelOption struct {
	MaxConcurrency int // # of goroutines to invoke at max. 0 means the same # of slice length
}

// DefaultParallelOption is a default ParallelOption value
var DefaultParallelOption = &ParallelOption{
	MaxConcurrency: 0,
}

// Parallel is spawn `fun` in parallel.
// `fun` must be type of `func(int, *T) error`, where list is []T or []*T).
func Parallel(list interface{}, option *ParallelOption, fun interface{}) error {
	if option == nil {
		option = DefaultParallelOption
	}
	a1 := reflect.ValueOf(list)
	f := reflect.ValueOf(fun)
	fType := f.Type()
	shouldUsePtr := a1.Type().Elem().Kind() == reflect.Struct

	assertSlice(a1)
	assertSliceFun(fType)
	if fType.NumOut() != 1 || !fType.Out(0).Implements(errorType) {
		panic(fmt.Errorf("SliceFuncError: the second function must return an error"))
	}

	l := a1.Len()
	if l == 0 {
		return nil
	}
	n := option.MaxConcurrency
	if l < n {
		n = l
	}
	if n <= 0 {
		n = 1
	}

	errors := SliceError(make([]error, l))
	eachSize := l / n
	a := reflect.ValueOf(SplitByLength(list, eachSize))

	var wg sync.WaitGroup

	for i := 0; i < a.Len(); i++ {
		wg.Add(1)
		v := a.Index(i) // still slice as we used SplitByEach
		go func(i int, v reflect.Value) {
			defer wg.Done()
			for j := 0; j < v.Len(); j++ {
				idx := i*eachSize + j
				func() {
					defer func() {
						if x := recover(); x != nil {
							errors[idx] = fmt.Errorf("%v", x)
						}
					}()
					v1 := v.Index(j)
					if shouldUsePtr {
						v1 = v1.Addr()
					}
					out := f.Call([]reflect.Value{reflect.ValueOf(idx), v1})[0]
					if !out.IsNil() {
						errors[idx] = out.Interface().(error)
					}
				}()
			}
		}(i, v)
	}
	wg.Wait()
	any := false
	for i := 0; i < l; i++ {
		if errors[i] != nil {
			any = true
		}
	}
	if any {
		return errors
	}
	return nil
}
