package datastore

import (
	"fmt"

	"context"

	"cloud.google.com/go/datastore"
	"github.com/yssk22/go/gcp"
)

// LoggerKey is a key for logger in this package
const LoggerKey = "gae.datastore"

// NewKey returns a new *datastore.Key for `kind`.
// if k is *datastore.Key, it returns the same object.
// if k is not a string nor an int, k is converted by fmt.Sprintf("%s").
func NewKey(ctx context.Context, kind string, k interface{}) *datastore.Key {
	var key *datastore.Key
	switch k.(type) {
	case string:
		key = datastore.NameKey(kind, k.(string), nil)
	case []byte:
		key = datastore.NameKey(kind, string(k.([]byte)), nil)
	case *datastore.Key:
		key = k.(*datastore.Key)
	default:
		key = datastore.NameKey(kind, fmt.Sprintf("%s", k), nil)
	}
	key.Namespace = gcp.CurrentNamespace(ctx)
	return key
}

// GetCacheKey returns a string representation for the cache key
func GetCacheKey(k *datastore.Key) string {
	return fmt.Sprintf("datastore.%s", k.Encode())
}

// IsDatastoreError returns true if err is not ErrNoSuchEntity
func IsDatastoreError(err error) bool {
	if err == nil {
		return false
	}

	if merror, ok := err.(datastore.MultiError); ok {
		for _, e := range merror {
			if e != nil && e != datastore.ErrNoSuchEntity {
				return true
			}
		}
		return false
	}
	return true
}

// NormalizeKeys to normalize keys from []string, []interface{} to []*datastore.Key
func NormalizeKeys(ctx context.Context, kind string, keys interface{}) ([]*datastore.Key, error) {
	var dsKeys []*datastore.Key
	switch t := keys.(type) {
	case []string:
		tmp := keys.([]string)
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = NewKey(ctx, kind, s)
		}
	case []interface{}:
		tmp := keys.([]interface{})
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = NewKey(ctx, kind, s)
		}
	case []*datastore.Key:
		dsKeys = keys.([]*datastore.Key)
	default:
		return nil, fmt.Errorf("unsupported keys type: %s", t)
	}
	return dsKeys, nil
}
