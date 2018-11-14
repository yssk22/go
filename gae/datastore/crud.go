package datastore

import (
	"context"
	"fmt"
	"reflect"

	"github.com/yssk22/go/iterator/slice"
	"github.com/yssk22/go/gae/memcache"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xlog"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var datastoreLoggerKey = struct{}{}

// CrudEntsLimit is a limit of the number of entities that can be handled in one put or delete transaction
const CrudEntsLimit = 200

var (
	// ErrTooManyEnts is returned when the user passes too many entities to PutMulti or DeleteMulti.
	ErrTooManyEnts = fmt.Errorf("too many entities to operate (max: %d)", CrudEntsLimit)
)

// GetMulti is wrapper for google.golang.org/appengine/datastore.GetMulti
func GetMulti(ctx context.Context, keys []*datastore.Key, entities interface{}, options ...CRUDOption) error {
	var err error
	var memKeys []string
	cfg := newCrudConfig(options...)
	if cfg.Namespace != nil {
		ctx, err = appengine.Namespace(ctx, *(cfg.Namespace))
		xerrors.MustNil(err)
	}
	size := len(keys)
	if size == 0 {
		return nil
	}
	if size > CrudEntsLimit {
		return ErrTooManyEnts
	}

	if !cfg.NoCache {
		memKeys = make([]string, size, size)
		for i := range keys {
			memKeys[i] = GetMemcacheKey(keys[i])
		}
		err = memcache.GetMulti(ctx, memKeys, entities)
		if err == nil {
			return nil
		}
	}
	// TODO: only fetch entities[i] that is nil
	cacheMissingKeys := make([]*datastore.Key, 0)
	cacheMissingEnts := slice.Filter(entities, func(i int, v interface{}) bool {
		if reflect.ValueOf(v).IsNil() {
			cacheMissingKeys = append(cacheMissingKeys, keys[i])
			return false
		}
		return true
	})

	// we check if err is an datastore error not to return "no such entity" error.
	if err = datastore.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts); IsDatastoreError(err) {
		_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace(), keys[0].Kind()), datastoreLoggerKey)
		logger.Fatalf("database error: %v", err)
		return err
	}
	var cacheKeyIndex = 0
	var entsFromDatastore = reflect.ValueOf(cacheMissingEnts)
	ents := reflect.ValueOf(entities)
	for i:=0; i<size; i++ {
		elm := ents.Index(i)
		if elm.IsNil() {
			vv := entsFromDatastore.Index(cacheKeyIndex)
			if !vv.IsNil() {
				elm.Set(vv)
			}
			cacheKeyIndex++		
		}
	}
	if !cfg.NoCache {
		if err := memcache.SetMulti(ctx, memKeys, entities); err != nil {
			_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace(), keys[0].Kind()), datastoreLoggerKey)
			logger.Warnf("could not update the datastore cache: %v", err)
		}		
	}
	return nil
}

// PutMulti is wrapper for google.golang.org/appengine/datastore.PutMulti
func PutMulti(ctx context.Context, keys []*datastore.Key, ent interface{}, options ...CRUDOption) ([]*datastore.Key, error) {
	var err error
	cfg := newCrudConfig(options...)
	if cfg.Namespace != nil {
		ctx, err = appengine.Namespace(ctx, *(cfg.Namespace))
		xerrors.MustNil(err)
	}
	size := len(keys)
	if size == 0 {
		return []*datastore.Key{}, nil
	}
	if size > CrudEntsLimit {
		return nil, ErrTooManyEnts
	}
	keys, err = datastore.PutMulti(ctx, keys, ent)
	if IsDatastoreError(err) {
		_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace(), keys[0].Kind()), datastoreLoggerKey)
		logger.Fatalf("database error: %v", err)
		return nil, err
	}

	if !cfg.NoCache {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = GetMemcacheKey(keys[i])
		}
		if err = memcache.DeleteMulti(ctx, memKeys); memcache.IsMemcacheError(err) {
			_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace(), keys[0].Kind()), datastoreLoggerKey)
			logger.Warnf("could not update the datastore cache: %v", err)
		}
	}
	return keys, nil
}

// DeleteMulti is wrapper for google.golang.org/appengine/datastore.DeleteMulti
func DeleteMulti(ctx context.Context, keys []*datastore.Key, options ...CRUDOption) error {
	var err error
	cfg := newCrudConfig(options...)
	if cfg.Namespace != nil {
		ctx, err = appengine.Namespace(ctx, *(cfg.Namespace))
		xerrors.MustNil(err)
	}
	size := len(keys)
	if size == 0 {
		return nil
	}
	if size > CrudEntsLimit {
		return ErrTooManyEnts
	}
	err = datastore.DeleteMulti(ctx, keys)
	if IsDatastoreError(err) {
		_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace(), keys[0].Kind()), datastoreLoggerKey)
		logger.Fatalf("database error: %v", err)
		return err
	}

	if !cfg.NoCache {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = GetMemcacheKey(keys[i])
		}
		if err = memcache.DeleteMulti(ctx, memKeys); memcache.IsMemcacheError(err) {
			_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace(), keys[0].Kind()), datastoreLoggerKey)
			logger.Warnf("could not update the datastore cache: %v", err)
		}
	}
	return nil
}

type crudConfig struct {
	NoCache              bool
	Namespace            *string
}

func newCrudConfig(options ...CRUDOption) *crudConfig {
	opts := &crudConfig{}
	for _, f := range options {
		opts = f(opts)
	}
	return opts
}

// CRUDOption is a function to configure CRUD operation
type CRUDOption func(*crudConfig) *crudConfig

// DontCache to tell do not cache.
func DontCache() CRUDOption {
	return CRUDOption(func(opts *crudConfig) *crudConfig {
		opts.NoCache = true
		return opts
	})
}

// Namespace to set the namespace
func Namespace(ns string) CRUDOption {
	return CRUDOption(func(opts *crudConfig) *crudConfig {
		opts.Namespace = &ns
		return opts
	})
}
