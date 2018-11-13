package datastore

import (
	"fmt"

	"context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// LoggerKey is a key for logger in this package
const LoggerKey = "gae.datastore"

// MaxEntitiesPerUpdate is a limit of the number of entities that can be handled in one put or delete transaction
const MaxEntitiesPerUpdate = 200

var (
	// ErrTooManyEnts is returned when the user passes too many entities to PutMulti or DeleteMulti.
	ErrTooManyEnts = fmt.Errorf("ent: too many documents given to put or delete (max is %d)", MaxEntitiesPerUpdate)
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

// GetMulti is wrapper for google.golang.org/appengine/datastore.GetMulti
// to support +1000 keys
func GetMulti(ctx context.Context, keys []*datastore.Key, ent interface{}) error {
	// TODO: support +1000 keys
	return datastore.GetMulti(ctx, keys, ent)
}

// PutMulti is wrapper for google.golang.org/appengine/datastore.PutMulti
// to support +1000 keys
func PutMulti(ctx context.Context, keys []*datastore.Key, ent interface{}) ([]*datastore.Key, error) {
	// TODO: support +1000 keys
	return datastore.PutMulti(ctx, keys, ent)
}

// DeleteMulti is wrapper for google.golang.org/appengine/datastore.DeleteMulti
// to support +1000 keys
func DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	// TODO: support +1000 keys
	return datastore.DeleteMulti(ctx, keys)
}

// CRUDOption to represent options to crud operation datastore
type CRUDOption struct {
	NoCache              bool
	NoTimestampUpdate    bool
	NoSearchIndexing     bool
	IgnoreSearchIndexing bool
	Namespace            *string
}

// Option is a function to configure *Option
type Option func(*CRUDOption) *CRUDOption

// DontCache to tell do not cache.
func DontCache() Option {
	return Option(func(opts *CRUDOption) *CRUDOption {
		opts.NoCache = true
		return opts
	})
}

// DontUpdateTimestamp to tell do not update the timestamp.
func DontUpdateTimestamp() Option {
	return Option(func(opts *CRUDOption) *CRUDOption {
		opts.NoTimestampUpdate = true
		return opts
	})
}

// DontIndexForSearch to tell do not update the search index.
func DontIndexForSearch() Option {
	return Option(func(opts *CRUDOption) *CRUDOption {
		opts.NoSearchIndexing = true
		return opts
	})
}

// IgnoreSearchIndexError to tell ignore errors by search indexing.
func IgnoreSearchIndexError() Option {
	return Option(func(opts *CRUDOption) *CRUDOption {
		opts.IgnoreSearchIndexing = true
		return opts
	})
}

// Namespace to set the namespace
func Namespace(ns string) Option {
	return Option(func(opts *CRUDOption) *CRUDOption {
		opts.Namespace = &ns
		return opts
	})
}

// NewCRUDOption returns a new *CRUDOption instance
func NewCRUDOption(options ...Option) *CRUDOption {
	opts := &CRUDOption{}
	for _, f := range options {
		f(opts)
	}
	return opts
}
