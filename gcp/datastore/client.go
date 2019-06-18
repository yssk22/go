package datastore

import (
	"context"
	"fmt"
	"reflect"

	"cloud.google.com/go/datastore"
	"github.com/yssk22/go/cache"
	"github.com/yssk22/go/iterator/slice"
	"github.com/yssk22/go/x/xcontext"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xlog"
)

// Client is a wrapper for datastore.Client
type Client struct {
	inner  *datastore.Client
	config *clientConfig
}

var contextClientKey = xcontext.NewKey("client")

// WithClient setup a *Client for the current context which can be referred by FromContext
func WithClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, contextClientKey, client)
}

// FromContext returns a *Client for the current context
func FromContext(ctx context.Context) *Client {
	return ctx.Value(contextClientKey).(*Client)
}

// NewClient returns a new Client
func NewClient(ctx context.Context, projectID string, options ...Option) *Client {
	config := newClientConfig(options...)
	inner, err := datastore.NewClient(ctx, projectID)
	xerrors.MustNil(err)
	return &Client{
		inner:  inner,
		config: config,
	}
}

// NewClientFromClient returns a new *Client from the *datastore.Client
func NewClientFromClient(ctx context.Context, c *datastore.Client, options ...Option) *Client {
	return &Client{
		inner:  c,
		config: newClientConfig(options...),
	}
}

type clientConfig struct {
	Cache     cache.Cache
	Namespace *string
}

func newClientConfig(options ...Option) *clientConfig {
	opts := &clientConfig{}
	for _, f := range options {
		opts = f(opts)
	}
	return opts
}

// Option is a function to configure CRUD operation
type Option func(*clientConfig) *clientConfig

// Cache to set the cache storage
func Cache(c cache.Cache) Option {
	return Option(func(opts *clientConfig) *clientConfig {
		opts.Cache = c
		return opts
	})
}

// Namespace to set the namespace
func Namespace(ns string) Option {
	return Option(func(opts *clientConfig) *clientConfig {
		opts.Namespace = &ns
		return opts
	})
}

var datastoreLoggerKey = struct{}{}

// CrudEntsLimit is a limit of the number of entities that can be handled in one put or delete transaction
const CrudEntsLimit = 200

var (
	// ErrTooManyEnts is returned when the user passes too many entities to PutMulti or DeleteMulti.
	ErrTooManyEnts = fmt.Errorf("too many entities to operate (max: %d)", CrudEntsLimit)
)

// GetMulti is wrapper for google.golang.org/appengine/datastore.GetMulti
func (c *Client) GetMulti(ctx context.Context, keys []*datastore.Key, entities interface{}, options ...Option) error {
	var err error
	var memKeys []string
	size := len(keys)
	if size == 0 {
		return nil
	}
	if size > CrudEntsLimit {
		return ErrTooManyEnts
	}

	if c.config.Cache != nil {
		memKeys = make([]string, size, size)
		for i := range keys {
			memKeys[i] = GetCacheKey(keys[i])
		}
		err = c.config.Cache.GetMulti(ctx, memKeys, entities)
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
	if err = c.inner.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts); IsDatastoreError(err) {
		_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace, keys[0].Kind), datastoreLoggerKey)
		logger.Fatalf("database error: %v", err)
		return err
	}
	var cacheKeyIndex = 0
	var entsFromDatastore = reflect.ValueOf(cacheMissingEnts)
	ents := reflect.ValueOf(entities)
	for i := 0; i < size; i++ {
		elm := ents.Index(i)
		if elm.IsNil() {
			vv := entsFromDatastore.Index(cacheKeyIndex)
			if !vv.IsNil() {
				elm.Set(vv)
			}
			cacheKeyIndex++
		}
	}
	if c.config.Cache != nil {
		if err := c.config.Cache.SetMulti(ctx, memKeys, entities); err != nil {
			_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace, keys[0].Kind), datastoreLoggerKey)
			logger.Warnf("could not update the datastore cache: %v", err)
		}
	}
	return nil
}

// PutMulti is wrapper for google.golang.org/appengine/datastore.PutMulti
func (c *Client) PutMulti(ctx context.Context, keys []*datastore.Key, ent interface{}, options ...Option) ([]*datastore.Key, error) {
	var err error
	size := len(keys)
	if size == 0 {
		return []*datastore.Key{}, nil
	}
	if size > CrudEntsLimit {
		return nil, ErrTooManyEnts
	}
	_, err = c.inner.PutMulti(ctx, keys, ent)
	if IsDatastoreError(err) {
		_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace, keys[0].Kind), datastoreLoggerKey)
		logger.Fatalf("database error: %v", err)
		return nil, err
	}

	if c.config.Cache != nil {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = GetCacheKey(keys[i])
		}
		if err = c.config.Cache.DeleteMulti(ctx, memKeys); err != nil {
			_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace, keys[0].Kind), datastoreLoggerKey)
			logger.Warnf("could not update the datastore cache: %v", err)
		}
	}
	return keys, nil
}

// DeleteMulti is wrapper for google.golang.org/appengine/datastore.DeleteMulti
func (c *Client) DeleteMulti(ctx context.Context, keys []*datastore.Key, options ...Option) error {
	var err error
	size := len(keys)
	if size == 0 {
		return nil
	}
	if size > CrudEntsLimit {
		return ErrTooManyEnts
	}
	err = c.inner.DeleteMulti(ctx, keys)
	if IsDatastoreError(err) {
		_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace, keys[0].Kind), datastoreLoggerKey)
		logger.Fatalf("database error: %v", err)
		return err
	}

	if c.config.Cache != nil {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = GetCacheKey(keys[i])
		}
		if err = c.config.Cache.DeleteMulti(ctx, memKeys); err != nil {
			_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("datastore.%s.%s", keys[0].Namespace, keys[0].Kind), datastoreLoggerKey)
			logger.Warnf("could not update the datastore cache: %v", err)
		}
	}
	return nil
}

// GetAll fills the query result into dst and returns corresponding *datastore.Key
func (c *Client) GetAll(ctx context.Context, q *Query, dst interface{}) ([]*datastore.Key, error) {
	return c.inner.GetAll(ctx, q.inner, dst)
}

// Run runs a query and returns *datastore.Iterator
func (c *Client) Run(ctx context.Context, q *Query) (*datastore.Iterator, error) {
	return c.inner.Run(ctx, q.inner), nil
}

// Count returns a count
func (c *Client) Count(ctx context.Context, q *Query) (int, error) {
	return c.inner.Count(ctx, q.inner)
}

// DeleteAll deletes the all `kind` entities stored in datastore
func (c *Client) DeleteAll(ctx context.Context, kind string) error {
	const batchSize = 300
	var keys []*datastore.Key
	var dummy []interface{}
	var err error
	for {
		if keys, err = c.inner.GetAll(ctx, datastore.NewQuery(kind).KeysOnly().Limit(batchSize), dummy); err != nil {
			return fmt.Errorf("delete_all: error retrieving keys: %v", err)
		}
		if err := c.inner.DeleteMulti(ctx, keys); err != nil {
			return fmt.Errorf("delete_all: error deleting keys: %v", err)
		}
		count, err := c.inner.Count(ctx, datastore.NewQuery(kind).KeysOnly())
		if err != nil {
			return fmt.Errorf("delete_all: error checking remaining keys: %v", err)
		}
		if count == 0 {
			break
		}
	}
	return nil
}
