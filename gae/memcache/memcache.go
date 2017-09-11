// Package memcache provides some utilities for memcache access
package memcache

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/speedland/go/web/response"

	"context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

// Set sets the value on memcache with key
func Set(ctx context.Context, key string, value interface{}) error {
	return SetWithExpire(ctx, key, value, time.Duration(0))
}

// SetWithExpire is like Set and with expiration
func SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	err := SetMultiWithExpire(ctx, []string{key}, []interface{}{value}, expire)
	if e, ok := err.(appengine.MultiError); ok {
		return e[0]
	}
	return nil
}

// SetMulti is to support multiple keys. `value`` must be []Kind or []*Kind
func SetMulti(ctx context.Context, keys []string, values interface{}) error {
	return SetMultiWithExpire(ctx, keys, values, time.Duration(0))
}

// SetMultiWithExpire is like SetMulti and with expiration
func SetMultiWithExpire(ctx context.Context, keys []string, values interface{}, expire time.Duration) error {
	items := make([]*memcache.Item, len(keys), len(keys))
	srcValue := reflect.ValueOf(values)
	for i, key := range keys {
		obj := srcValue.Index(i).Interface()
		buff, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		items[i] = &memcache.Item{
			Key:        key,
			Value:      buff,
			Expiration: expire,
		}
	}
	if len(keys) == 1 {
		return memcache.Set(ctx, items[0])
	}
	return memcache.SetMulti(ctx, items)
}

const jsonNullValue = "null"

// GetMulti gets the values by multiple keys
func GetMulti(ctx context.Context, keys []string, dst interface{}) error {
	if len(keys) == 0 {
		return nil
	}
	var errors = make([]error, len(keys), len(keys))
	var numErrors = 0
	// GetMulti reutrns map[string]*memcache.Item
	hits, err := memcache.GetMulti(ctx, keys)
	if err != nil {
		return err
	}
	if err == nil {
		dstValue := reflect.ValueOf(dst)
		var isPtr = dstValue.Index(0).Kind() == reflect.Ptr
		for i, key := range keys {
			item, ok := hits[key]
			if !ok {
				errors[i] = memcache.ErrCacheMiss
				numErrors++
				continue
			}
			objValue := dstValue.Index(i)
			if isPtr {
				if string(item.Value) == jsonNullValue {
					continue
				}
				obj := reflect.New(objValue.Type().Elem()).Interface()
				if err = json.Unmarshal(item.Value, obj); err != nil {
					errors[i] = err
					numErrors++
					continue
				}
				objValue.Set(reflect.ValueOf(obj))
			} else {
				obj := objValue.Interface()
				if err = json.Unmarshal(item.Value, &obj); err != nil {
					errors[i] = err
					numErrors++
				}
			}
		}
	}
	if numErrors > 0 {
		return appengine.MultiError(errors)
	}
	return nil
}

// Get gets the object into dst by `key`
func Get(ctx context.Context, key string, dst interface{}) error {
	err := GetMulti(ctx, []string{key}, []interface{}{dst})
	if e, ok := err.(appengine.MultiError); ok {
		return e[0]
	}
	return err
}

// Delete deletes the key
func Delete(ctx context.Context, key string) {
	memcache.Delete(ctx, key)
}

// DeleteMulti deletes the multiple keys
func DeleteMulti(ctx context.Context, keys []string) error {
	return memcache.DeleteMulti(ctx, keys)
}

// Exists returns if the key exists
func Exists(ctx context.Context, key string) bool {
	_, err := memcache.Get(ctx, key)
	return err == nil
}

// IsMemcacheError returns whether the error is by Memcache
func IsMemcacheError(e error) bool {
	if e == nil {
		return false
	}
	if e == memcache.ErrCacheMiss {
		return false
	}
	if merr, ok := e.(appengine.MultiError); ok {
		for _, e := range merr {
			if e != nil && e != memcache.ErrCacheMiss {
				return true
			}
		}
		return false
	}
	return true
}

func init() {
	RegisterCachableHeaderKeys(
		response.ContentType,
	)
}
