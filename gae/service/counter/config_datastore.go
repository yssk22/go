// Code generated by "ent -type=Config"; DO NOT EDIT

package counter

import (
	"context"
	"fmt"
	"github.com/speedland/go/ent"
	helper "github.com/speedland/go/gae/datastore"
	"github.com/speedland/go/gae/memcache"
	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xerrors"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func (c *Config) NewKey(ctx context.Context) *datastore.Key {
	return helper.NewKey(ctx, "CounterConfig", c.Key)
}

// UpdateByForm updates the fields by form values. All values should be validated
// before calling this function.
func (c *Config) UpdateByForm(form *keyvalue.GetProxy) {
}

// NewConfig returns a new *Config with default field values.
func NewConfig() *Config {
	c := &Config{}
	return c
}

type ConfigKind struct {
	useDefaultIfNil           bool
	noCache                   bool
	noSearchIndexing          bool
	ignoreSearchIndexingError bool
	noTimestampUpdate         bool
	enforceNamespace          bool
	namespace                 string
}

type ConfigKindReplacer interface {
	Replace(*Config, *Config) *Config
}

type ConfigKindReplacerFunc func(*Config, *Config) *Config

func (f ConfigKindReplacerFunc) Replace(ent1 *Config, ent2 *Config) *Config {
	return f(ent1, ent2)
}

// DefaultConfigKind is a default value of *ConfigKind
var DefaultConfigKind = &ConfigKind{}

// ConfigKindLoggerKey is a logger key name for the ent
const ConfigKindLoggerKey = "ent.counter_config"

// ConfigQueryLoggerKey is a logger key name for the ent
const ConfigQueryLoggerKey = "ent.query.counter_config"

// EnforceNamespace enforces namespace for Get/Put/Delete or not.
func (k *ConfigKind) EnforceNamespace(ns string, b bool) *ConfigKind {
	k.enforceNamespace = b
	k.namespace = ns
	return k
}

func (k *ConfigKind) UseDefaultIfNil(b bool) *ConfigKind {
	k.useDefaultIfNil = b
	return k
}

// Get gets the kind entity from datastore
func (k *ConfigKind) Get(ctx context.Context, key interface{}) (*datastore.Key, *Config, error) {
	keys, ents, err := k.GetMulti(ctx, []interface{}{key})
	if err != nil {
		return nil, nil, err
	}
	return keys[0], ents[0], nil
}

// MustGet is like Get but returns only values and panic if error happens.
func (k *ConfigKind) MustGet(ctx context.Context, key interface{}) *Config {
	_, v, err := k.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMulti do Get with multiple keys. keys must be []string, []*datastore.Key, or []interface{}
func (k *ConfigKind) GetMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, []*Config, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ConfigKindLoggerKey)
	var err error
	var dsKeys []*datastore.Key
	var memKeys []string
	var ents []*Config
	if k.enforceNamespace {
		ctx, err = appengine.Namespace(ctx, k.namespace)
		if err != nil {
			return nil, nil, xerrors.Wrap(err, "cannot enforce namespace")
		}
	}
	dsKeys, err = k.normMultiKeys(ctx, keys)
	if err != nil {
		return nil, nil, xerrors.Wrap(err, "cannot normalize keys")
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil, nil
	}
	ents = make([]*Config, size, size)
	// Memcache access
	if !k.noCache {
		logger.Debugf("Trying to get entities from memcache...")
		memKeys = make([]string, size, size)
		for i := range dsKeys {
			memKeys[i] = ent.GetMemcacheKey(dsKeys[i])
		}
		err = memcache.GetMulti(ctx, memKeys, ents)
		if err == nil {
			// Hit caches on all keys!!
			return dsKeys, ents, nil
		}
		logger.Debug(func(p *xlog.Printer) {
			p.Println("CounterConfig#GetMulti [Memcache]")
			for i := 0; i < size; i++ {
				s := fmt.Sprintf("%v", ents[i])
				if len(s) > 20 {
					p.Printf("\t%s - %s...\n", memKeys[i], s[:20])
				} else {
					p.Printf("\t%s - %s\n", memKeys[i], s)
				}
				if i >= 20 {
					p.Printf("\t...(and %d ents)\n", size-i)
					break
				}
			}
		})
	}

	key2Idx := make(map[*datastore.Key]int)
	cacheMissingKeys := make([]*datastore.Key, 0)
	for i := range ents {
		if ents[i] == nil {
			key2Idx[dsKeys[i]] = i
			cacheMissingKeys = append(cacheMissingKeys, dsKeys[i])
		}
	}
	cacheMissingSize := len(cacheMissingKeys)

	// Datastore access
	cacheMissingEnts := make([]*Config, cacheMissingSize, cacheMissingSize)
	err = helper.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts)
	if helper.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, nil, xerrors.Wrap(err, "datastore error")
	}

	if k.useDefaultIfNil {
		for i := 0; i < cacheMissingSize; i++ {
			if cacheMissingEnts[i] == nil {
				cacheMissingEnts[i] = NewConfig()
				cacheMissingEnts[i].Key = cacheMissingKeys[i].StringID() // TODO: Support non-string key as ID
			}
		}
	}

	// merge cacheMissingEnts with ents.
	for i := range cacheMissingKeys {
		entIdx := key2Idx[cacheMissingKeys[i]]
		ents[entIdx] = cacheMissingEnts[i]
	}

	// create a cache
	if !k.noCache {
		cacheEnts := make([]*Config, 0)
		cacheKeys := make([]string, 0)
		for i := range ents {
			if ents[i] != nil {
				cacheEnts = append(cacheEnts, ents[i])
				cacheKeys = append(cacheKeys, memKeys[i])
			}
		}
		if len(cacheEnts) > 0 {
			if err := memcache.SetMulti(ctx, cacheKeys, cacheEnts); err != nil {
				logger.Warnf("Failed to create CounterConfig) caches: %v", err)
			}
		}
	}

	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"CounterConfig#GetMulti [Datastore] (UseDefault: %t, NoCache: %t)\n",
			k.useDefaultIfNil, k.noCache,
		)
		for i := 0; i < size; i++ {
			s := fmt.Sprintf("%v", ents[i])
			if len(s) > 20 {
				p.Printf("\t%s - %s...\n", dsKeys[i], s[:20])
			} else {
				p.Printf("\t%s - %s\n", dsKeys[i], s)
			}
			if i >= 20 {
				p.Printf("\t...(and %d ents)\n", size-i)
				break
			}
		}
	})
	return dsKeys, ents, nil
}

// MustGetMulti is like GetMulti but returns only values and panic if error happens.
func (k *ConfigKind) MustGetMulti(ctx context.Context, keys interface{}) []*Config {
	_, v, err := k.GetMulti(ctx, keys)
	if err != nil {
		panic(err)
	}
	return v
}

// Put puts the entity to datastore.
func (k *ConfigKind) Put(ctx context.Context, ent *Config) (*datastore.Key, error) {
	keys, err := k.PutMulti(ctx, []*Config{
		ent,
	})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustPut is like Put and panic if an error occurrs.
func (k *ConfigKind) MustPut(ctx context.Context, ent *Config) *datastore.Key {
	keys, err := k.Put(ctx, ent)
	if err != nil {
		panic(err)
	}
	return keys
}

// PutMulti do Put with multiple keys
func (k *ConfigKind) PutMulti(ctx context.Context, ents []*Config) ([]*datastore.Key, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ConfigKindLoggerKey)
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	if size == 0 {
		return nil, nil
	}
	if size >= ent.MaxEntsPerPutDelete {
		return nil, ent.ErrTooManyEnts
	}
	if k.enforceNamespace {
		ctx, err = appengine.Namespace(ctx, k.namespace)
		if err != nil {
			return nil, xerrors.Wrap(err, "cannot enforce namespace")
		}
	}
	dsKeys = make([]*datastore.Key, size, size)
	for i := range ents {
		if e, ok := interface{}(ents[i]).(ent.BeforeSave); ok {
			if err := e.BeforeSave(ctx); err != nil {
				return nil, err
			}
		}
		dsKeys[i] = ents[i].NewKey(ctx)
	}

	if !k.noTimestampUpdate {
		for i := range ents {
			ents[i].UpdatedAt = xtime.Now()
		}
	}

	_, err = helper.PutMulti(ctx, dsKeys, ents)
	if helper.IsDatastoreError(err) {
		return nil, xerrors.Wrap(err, "datastore error")
	}

	if !k.noCache {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = ent.GetMemcacheKey(dsKeys[i])
		}
		err := memcache.DeleteMulti(ctx, memKeys)
		if memcache.IsMemcacheError(err) {
			logger.Warnf("Failed to invalidate memcache keys: %v", err)
		}
	}

	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"CounterConfig#PutMulti [Datastore] (NoCache: %t)\n",
			k.noCache,
		)
		for i := 0; i < size; i++ {
			s := fmt.Sprintf("%v", ents[i])
			if len(s) > 20 {
				p.Printf("\t%s - %s...\n", dsKeys[i], s[:20])
			} else {
				p.Printf("\t%s - %s\n", dsKeys[i], s)
			}
			if i >= 20 {
				p.Printf("\t...(and %d ents)\n", size-i)
				break
			}
		}
	})

	return dsKeys, nil
}

// MustPutMulti is like PutMulti but panic if an error occurs
func (k *ConfigKind) MustPutMulti(ctx context.Context, ents []*Config) []*datastore.Key {
	keys, err := k.PutMulti(ctx, ents)
	if err != nil {
		panic(err)
	}
	return keys
}

func (k *ConfigKind) Replace(ctx context.Context, ent *Config, replacer ConfigKindReplacer) (*datastore.Key, *Config, error) {
	keys, ents, err := k.ReplaceMulti(ctx, []*Config{
		ent,
	}, replacer)
	if err != nil {
		return nil, ents[0], err
	}
	return keys[0], ents[0], err
}

func (k *ConfigKind) MustReplace(ctx context.Context, ent *Config, replacer ConfigKindReplacer) (*datastore.Key, *Config) {
	key, ent, err := k.Replace(ctx, ent, replacer)
	if err != nil {
		panic(err)
	}
	return key, ent
}

func (k *ConfigKind) ReplaceMulti(ctx context.Context, ents []*Config, replacer ConfigKindReplacer) ([]*datastore.Key, []*Config, error) {
	var size = len(ents)
	var dsKeys = make([]*datastore.Key, size, size)
	if size == 0 {
		return dsKeys, ents, nil
	}
	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
	}
	_, existing, err := k.GetMulti(ctx, dsKeys)
	if err != nil {
		return nil, ents, err
	}
	for i, exist := range existing {
		if exist != nil {
			ents[i] = replacer.Replace(exist, ents[i])
		}
	}
	_, err = k.PutMulti(ctx, ents)
	return dsKeys, ents, err
}

func (k *ConfigKind) MustReplaceMulti(ctx context.Context, ents []*Config, replacer ConfigKindReplacer) ([]*datastore.Key, []*Config) {
	keys, ents, err := k.ReplaceMulti(ctx, ents, replacer)
	if err != nil {
		panic(err)
	}
	return keys, ents
}

// Delete deletes the entity from datastore
func (k *ConfigKind) Delete(ctx context.Context, key interface{}) (*datastore.Key, error) {
	keys, err := k.DeleteMulti(ctx, []interface{}{key})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustDelete is like Delete but panic if an error occurs
func (k *ConfigKind) MustDelete(ctx context.Context, key interface{}) *datastore.Key {
	_key, err := k.Delete(ctx, key)
	if err != nil {
		panic(err)
	}
	return _key
}

// DeleteMulti do Delete with multiple keys
func (k *ConfigKind) DeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ConfigKindLoggerKey)
	var err error
	var dsKeys []*datastore.Key
	if k.enforceNamespace {
		ctx, err = appengine.Namespace(ctx, k.namespace)
		if err != nil {
			return nil, xerrors.Wrap(err, "cannot enforce namespace")
		}
	}
	dsKeys, err = k.normMultiKeys(ctx, keys)
	if err != nil {
		return nil, err
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil
	}
	if size >= ent.MaxEntsPerPutDelete {
		return nil, ent.ErrTooManyEnts
	}

	// Datastore access
	err = helper.DeleteMulti(ctx, dsKeys)
	if helper.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, xerrors.Wrap(err, "datastore error")
	}

	if !k.noCache {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = ent.GetMemcacheKey(dsKeys[i])
		}
		err = memcache.DeleteMulti(ctx, memKeys)
		if memcache.IsMemcacheError(err) {
			logger.Warnf("Failed to invalidate memcache keys: %v", err)
		}
	}

	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"CounterConfig#DeleteMulti [Datastore] (NoCache: %t)\n",
			k.noCache,
		)
		for i := 0; i < size; i++ {
			p.Printf("\t%s\n", dsKeys[i])
			if i >= 20 {
				p.Printf("\t...(and %d ents)\n", size-i)
				break
			}
		}
	})
	return dsKeys, nil
}

// MustDeleteMulti is like DeleteMulti but panic if an error occurs
func (k *ConfigKind) MustDeleteMulti(ctx context.Context, keys interface{}) []*datastore.Key {
	_keys, err := k.DeleteMulti(ctx, keys)
	if err != nil {
		panic(err)
	}
	return _keys
}

func (k *ConfigKind) normMultiKeys(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	var dsKeys []*datastore.Key
	switch t := keys.(type) {
	case []string:
		tmp := keys.([]string)
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = helper.NewKey(ctx, "CounterConfig", s)
		}
	case []interface{}:
		tmp := keys.([]interface{})
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = helper.NewKey(ctx, "CounterConfig", s)
		}
	case []*datastore.Key:
		dsKeys = keys.([]*datastore.Key)
	default:
		return nil, fmt.Errorf("unsupported keys type: %s", t)
	}
	return dsKeys, nil
}

// ConfigQuery helps to build and execute a query
type ConfigQuery struct {
	q       *helper.Query
	viaKeys *ConfigKind
}

func NewConfigQuery() *ConfigQuery {
	return &ConfigQuery{
		q: helper.NewQuery("CounterConfig", ConfigQueryLoggerKey),
	}
}

// Ancestor sets the ancestor filter
func (q *ConfigQuery) Ancestor(a lazy.Value) *ConfigQuery {
	q.q = q.q.Ancestor(a)
	return q
}

// Eq sets the "=" filter on the name field.
func (q *ConfigQuery) Eq(name string, value lazy.Value) *ConfigQuery {
	q.q = q.q.Eq(name, value)
	return q
}

// Lt sets the "<" filter on the "name" field.
func (q *ConfigQuery) Lt(name string, value lazy.Value) *ConfigQuery {
	q.q = q.q.Lt(name, value)
	return q
}

// Le sets the "<=" filter on the "name" field.
func (q *ConfigQuery) Le(name string, value lazy.Value) *ConfigQuery {
	q.q = q.q.Le(name, value)
	return q
}

// Gt sets the ">" filter on the "name" field.
func (q *ConfigQuery) Gt(name string, value lazy.Value) *ConfigQuery {
	q.q = q.q.Gt(name, value)
	return q
}

// Ge sets the ">=" filter on the "name" field.
func (q *ConfigQuery) Ge(name string, value lazy.Value) *ConfigQuery {
	q.q = q.q.Ge(name, value)
	return q
}

// Ne sets the "!=" filter on the "name" field.
func (q *ConfigQuery) Ne(name string, value lazy.Value) *ConfigQuery {
	q.q = q.q.Ne(name, value)
	return q
}

// Asc specifies ascending order on the given filed.
func (q *ConfigQuery) Asc(name string) *ConfigQuery {
	q.q = q.q.Asc(name)
	return q
}

// Desc specifies descending order on the given filed.
func (q *ConfigQuery) Desc(name string) *ConfigQuery {
	q.q = q.q.Desc(name)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ConfigQuery) Limit(n lazy.Value) *ConfigQuery {
	q.q = q.q.Limit(n)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ConfigQuery) Start(value lazy.Value) *ConfigQuery {
	q.q = q.q.Start(value)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ConfigQuery) End(value lazy.Value) *ConfigQuery {
	q.q = q.q.End(value)
	return q
}

// ViaKeys optimize to execute keys-only query then call k.GetMulti() to fetch values.
// This would reduce the datastore query and maximize the memcache usage if the query called many times in a short time window.
func (q *ConfigQuery) ViaKeys(k *ConfigKind) *ConfigQuery {
	q.viaKeys = k
	if k == nil {
		q.q = q.q.KeysOnly(false)
	} else {
		q.q = q.q.KeysOnly(true)
	}
	return q
}

// GetAll returns all key and value of the query.
func (q *ConfigQuery) GetAll(ctx context.Context) ([]*datastore.Key, []*Config, error) {
	if q.viaKeys != nil {
		keys, err := q.q.GetAll(ctx, nil)
		if err != nil {
			return nil, nil, err
		}
		_, ents, err := q.viaKeys.GetMulti(ctx, keys)
		return keys, ents, err
	}
	var v []*Config
	keys, err := q.q.KeysOnly(false).GetAll(ctx, &v)
	if err != nil {
		return nil, nil, err
	}
	return keys, v, err
}

// MustGetAll is like GetAll but panic if an error occurrs.
func (q *ConfigQuery) MustGetAll(ctx context.Context) ([]*datastore.Key, []*Config) {
	keys, values, err := q.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	return keys, values
}

// GetAllValues is like GetAll but returns only values
func (q *ConfigQuery) GetAllValues(ctx context.Context) ([]*Config, error) {
	_, v, err := q.GetAll(ctx)
	return v, err
}

// MustGetAllValues is like GetAllValues but panic if an error occurrs
func (q *ConfigQuery) MustGetAllValues(ctx context.Context) []*Config {
	_, v, err := q.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Count returns the count of entities
func (q *ConfigQuery) Count(ctx context.Context) (int, error) {
	return q.q.Count(ctx)
}

// MustCount returns the count of entities
func (q *ConfigQuery) MustCount(ctx context.Context) int {
	c, err := q.Count(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

type ConfigPagination struct {
	Start string           `json:"start"`
	End   string           `json:"end"`
	Count int              `json:"count,omitempty"`
	Data  []*Config        `json:"data"`
	Keys  []*datastore.Key `json:"-"`
}

// Run returns the a result as *ConfigPagination object
func (q *ConfigQuery) Run(ctx context.Context) (*ConfigPagination, error) {
	iter, err := q.q.Run(ctx)
	if err != nil {
		return nil, err
	}
	pagination := &ConfigPagination{}
	keys := []*datastore.Key{}
	data := []*Config{}
	start, err := iter.Cursor()
	if err != nil {
		return nil, fmt.Errorf("couldn't get the start cursor: %v", err)
	}
	pagination.Start = start.String()
	for {
		var key *datastore.Key
		var err error
		var ent Config
		if q.viaKeys == nil {
			key, err = iter.Next(&ent)
		} else {
			key, err = iter.Next(nil)
		}
		if err == datastore.Done {
			end, err := iter.Cursor()
			if err != nil {
				return nil, fmt.Errorf("couldn't get the end cursor: %v", err)
			}
			pagination.Keys = keys
			pagination.End = end.String()
			if q.viaKeys != nil {
				_, ents, err := q.viaKeys.GetMulti(ctx, keys)
				if err != nil {
					return nil, err
				}
				pagination.Data = ents
			} else {
				pagination.Data = data
			}
			return pagination, nil
		}
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		if q.viaKeys == nil {
			data = append(data, &ent)
		}
	}
}

// MustRun is like Run but panic if an error occurrs
func (q *ConfigQuery) MustRun(ctx context.Context) *ConfigPagination {
	p, err := q.Run(ctx)
	if err != nil {
		panic(err)
	}
	return p
}

// DeleteMatched deletes the all ents that match with the query.
// This func modify Limit/StartKey condition in the query so that you should restore it
// if you want to reuse the query.
func (k *ConfigKind) DeleteMatched(ctx context.Context, q *ConfigQuery) (int, error) {
	var numDeletes int
	var startKey string
	q.Limit(lazy.New(ent.MaxEntsPerPutDelete - 5))
	// TODO: canceling the context
	for {
		if startKey != "" {
			q.Start(lazy.New(startKey))
		}
		page := q.MustRun(ctx)
		if len(page.Keys) == 0 {
			return numDeletes, nil
		}
		_, err := k.DeleteMulti(ctx, page.Keys)
		if err != nil {
			return numDeletes, fmt.Errorf("couldn't delete matched ents: %v", err)
		} else {
			numDeletes += len(page.Keys)
		}
		startKey = page.End
	}
}
