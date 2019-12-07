package slice

import (
	"fmt"
	"reflect"
	"sync"
)

type ParallelConfig struct {
	maxConcurrency int
}

type ParallelOption func(*ParallelConfig) *ParallelConfig

// MaxConcurrency to configure max cncurrurency
func MaxConcurrency(n int) ParallelOption {
	return func(p *ParallelConfig) *ParallelConfig {
		p.maxConcurrency = n
		return p
	}
}

// Parallel is spawn `fun` in parallel.
// `fun` must be type of `func(int, *T) error`, where list is []T or []*T).
func Parallel(list interface{}, fun interface{}, options ...ParallelOption) error {
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
	var config = &ParallelConfig{
		maxConcurrency: l,
	}
	for _, opt := range options {
		config = opt(config)
	}
	n := config.maxConcurrency
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
