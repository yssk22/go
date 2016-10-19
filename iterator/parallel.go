package iterator

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

// ParallelSlice is like ParallelSlice but spawn goroutines up to `n` concurrency
func ParallelSlice(list interface{}, option *ParallelOption, fun interface{}) error {
	if option == nil {
		option = DefaultParallelOption
	}
	a1 := reflect.ValueOf(list)
	f := reflect.ValueOf(fun)
	assertSlice(a1)
	assertSliceFun(f)
	l := a1.Len()
	n := option.MaxConcurrency
	if n <= 0 || l < n {
		n = l
	}

	errors := SliceError(make([]error, l))
	eachSize := l / n
	a := reflect.ValueOf(SplitSliceByLength(list, eachSize))

	var wg sync.WaitGroup

	shouldUsePtr := a.Type().Elem().Kind() == reflect.Struct
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
