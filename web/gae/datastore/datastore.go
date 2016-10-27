package datastore

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// NewKey returns a new *datastore.Key for `kind`.
// if k is *datastore.Key, it returns the same object.
// if k is not a string nor an int, k is converted by fmt.Sprintf("%s").
func NewKey(ctx context.Context, kind string, k interface{}) *datastore.Key {
	switch k.(type) {
	case string:
		return datastore.NewKey(ctx, kind, k.(string), 0, nil)
	case []byte:
		return datastore.NewKey(ctx, kind, string(k.([]byte)), 0, nil)
	case int:
		return datastore.NewKey(ctx, kind, "", int64(k.(int)), nil)
	case int8:
		return datastore.NewKey(ctx, kind, "", int64(k.(int8)), nil)
	case int16:
		return datastore.NewKey(ctx, kind, "", int64(k.(int16)), nil)
	case int32:
		return datastore.NewKey(ctx, kind, "", int64(k.(int32)), nil)
	case int64:
		return datastore.NewKey(ctx, kind, "", k.(int64), nil)
	case *datastore.Key:
		return k.(*datastore.Key)
	default:
		return datastore.NewKey(ctx, kind, fmt.Sprintf("%s", k), 0, nil)
	}
}

// IsDatastoreError returns true if err is not ErrNoSuchEntity
func IsDatastoreError(err error) bool {
	if err == nil {
		return false
	}
	if merror, ok := err.(appengine.MultiError); ok {
		for _, e := range merror {
			if e != nil && e != datastore.ErrNoSuchEntity {
				return true
			}
		}
		return false
	}
	return true
}

func GetMulti(ctx context.Context, keys []*datastore.Key, ent interface{}) error {
	// size := len(keys)
	// if size >= 1000 {
	// 	keyClusters := slice.SplitByLength(keys, 999).([][]*datastore.Key)
	// 	entClusters := slice.SplitByLength(ent, 999)
	// }
	// TODO: support +1000 keys
	return datastore.GetMulti(ctx, keys, ent)
}
