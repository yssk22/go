package cache

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/yssk22/go/iterator/slice"
	"github.com/yssk22/go/x/xerrors"
)

// MemoryCache is an example type that implements Cache in a single process environment.
type MemoryCache struct {
	m sync.Map
}

// SetMulti implements Cache#GetMulti
func (mc *MemoryCache) SetMulti(ctx context.Context, keys []string, values interface{}) error {
	return slice.ForEach(values, func(i int, v interface{}) error {
		mc.m.Store(keys[i], v)
		return nil
	})
}

var (
	ErrInvalidDstType   = errors.New("datastore: dst has invalid type")
	ErrInvalidDstLength = errors.New("datastore: key and dst slices have different length")
)

// GetMulti implements Cache#GetMulti
func (mc *MemoryCache) GetMulti(ctx context.Context, keys []string, dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Slice {
		return ErrInvalidDstType
	}
	if v.Len() != len(keys) {
		return ErrInvalidDstLength
	}
	errors := xerrors.NewMultiError(len(keys))
	slice.ForEach(keys, func(i int, k string) error {
		value, ok := mc.m.Load(k)
		if ok {
			vdst := v.Index(i)
			vsrc := reflect.ValueOf(value)
			vdstType := vdst.Type()
			vsrcType := vsrc.Type()
			if vdstType == vsrcType {
				vdst.Set(vsrc)
				errors = append(errors, nil)
			} else {
				if vdstType.Kind() == reflect.Ptr {
					// vdst: *A, vsrc: A
					if vdstType == reflect.PtrTo(vsrcType) {
						n := reflect.New(vsrcType)
						n.Elem().Set(vsrc)
						vdst.Set(n)
					} else {
						errors = append(errors, ErrInvalidDstType)
					}
				} else {
					// vdst: A, vsrc: *A
					vdstAddr := vdst.Addr()
					if vdstAddr.Type() == vsrcType {
						vdst.Set(vsrc.Elem())
						errors = append(errors, nil)
					} else {
						errors = append(errors, ErrInvalidDstType)
					}
				}
			}
		} else {
			errors = append(errors, ErrCacheKeyNotFound(keys[i]))
		}
		return nil
	})
	if errors.HasError() {
		return errors
	}
	return nil
}

// DeleteMulti implements Cache#DeleteMulti
func (mc *MemoryCache) DeleteMulti(ctx context.Context, keys []string) error {
	for k, _ := range keys {
		mc.m.Delete(k)
	}
	return nil
}
