// Code generated by "ent -type=Shard"; DO NOT EDIT

package counter

import (
	"fmt"
	"github.com/speedland/go/ent"
	helper "github.com/speedland/go/gae/datastore"
	"github.com/speedland/go/gae/memcache"
	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func (s *Shard) NewKey(ctx context.Context) *datastore.Key {
	return helper.NewKey(ctx, "CounterShared", s.Key)
}

// UpdateByForm updates the fields by form values. All values should be validated
// before calling this function.
func (s *Shard) UpdateByForm(form *keyvalue.GetProxy) {
}

// NewShard returns a new *Shard with default field values.
func NewShard() *Shard {
	s := &Shard{}
	return s
}

type ShardKind struct {
	BeforeSave                func(ent *Shard) error
	AfterSave                 func(ent *Shard) error
	useDefaultIfNil           bool
	noCache                   bool
	noSearchIndexing          bool
	ignoreSearchIndexingError bool
	noTimestampUpdate         bool
}

// DefaultShardKind is a default value of *ShardKind
var DefaultShardKind = &ShardKind{}

// ShardKindLoggerKey is a logger key name for the ent
const ShardKindLoggerKey = "ent.counter_shared"

func (k *ShardKind) UseDefaultIfNil(b bool) *ShardKind {
	k.useDefaultIfNil = b
	return k
}

// Get gets the kind entity from datastore
func (k *ShardKind) Get(ctx context.Context, key interface{}) (*datastore.Key, *Shard, error) {
	keys, ents, err := k.GetMulti(ctx, []interface{}{key})
	if err != nil {
		return nil, nil, err
	}
	return keys[0], ents[0], nil
}

// MustGet is like Get but returns only values and panic if error happens.
func (k *ShardKind) MustGet(ctx context.Context, key interface{}) *Shard {
	_, v, err := k.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMulti do Get with multiple keys. keys must be []string, []*datastore.Key, or []interface{}
func (k *ShardKind) GetMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, []*Shard, error) {
	var logger = xlog.WithContext(ctx).WithKey(ShardKindLoggerKey)
	var dsKeys, err = k.normMultiKeys(ctx, keys)
	if err != nil {
		return nil, nil, err
	}
	var size = len(dsKeys)
	var memKeys []string
	var ents []*Shard
	if size == 0 {
		return nil, nil, nil
	}
	ents = make([]*Shard, size, size)
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
			p.Println("CounterShared#GetMulti [Memcache]")
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
	cacheMissingEnts := make([]*Shard, cacheMissingSize, cacheMissingSize)
	err = helper.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts)
	if helper.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, nil, err
	}

	if k.useDefaultIfNil {
		for i := 0; i < cacheMissingSize; i++ {
			if cacheMissingEnts[i] == nil {
				cacheMissingEnts[i] = NewShard()
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
		cacheEnts := make([]*Shard, 0)
		cacheKeys := make([]string, 0)
		for i := range ents {
			if ents[i] != nil {
				cacheEnts = append(cacheEnts, ents[i])
				cacheKeys = append(cacheKeys, memKeys[i])
			}
		}
		if len(cacheEnts) > 0 {
			if err := memcache.SetMulti(ctx, cacheKeys, cacheEnts); err != nil {
				logger.Warnf("Failed to create CounterShared) caches: %v", err)
			}
		}
	}

	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"CounterShared#GetMulti [Datastore] (UseDefault: %t, NoCache: %t)\n",
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
func (k *ShardKind) MustGetMulti(ctx context.Context, keys interface{}) []*Shard {
	_, v, err := k.GetMulti(ctx, keys)
	if err != nil {
		panic(err)
	}
	return v
}

// Put puts the entity to datastore.
func (k *ShardKind) Put(ctx context.Context, ent *Shard) (*datastore.Key, error) {
	keys, err := k.PutMulti(ctx, []*Shard{
		ent,
	})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustPut is like Put and panic if an error occurrs.
func (k *ShardKind) MustPut(ctx context.Context, ent *Shard) *datastore.Key {
	keys, err := k.Put(ctx, ent)
	if err != nil {
		panic(err)
	}
	return keys
}

// PutMulti do Put with multiple keys
func (k *ShardKind) PutMulti(ctx context.Context, ents []*Shard) ([]*datastore.Key, error) {
	var size = len(ents)
	var dsKeys []*datastore.Key
	if size == 0 {
		return nil, nil
	}
	if size >= ent.MaxEntsPerPutDelete {
		return nil, ent.ErrTooManyEnts
	}
	logger := xlog.WithContext(ctx).WithKey(ShardKindLoggerKey)

	dsKeys = make([]*datastore.Key, size, size)
	for i := range ents {
		if k.BeforeSave != nil {
			if err := k.BeforeSave(ents[i]); err != nil {
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

	_, err := helper.PutMulti(ctx, dsKeys, ents)
	if helper.IsDatastoreError(err) {
		return nil, err
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
			"CounterShared#PutMulti [Datastore] (NoCache: %t)\n",
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
func (k *ShardKind) MustPutMulti(ctx context.Context, ents []*Shard) []*datastore.Key {
	keys, err := k.PutMulti(ctx, ents)
	if err != nil {
		panic(err)
	}
	return keys
}

// Delete deletes the entity from datastore
func (k *ShardKind) Delete(ctx context.Context, key interface{}) (*datastore.Key, error) {
	keys, err := k.DeleteMulti(ctx, []interface{}{key})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustDelete is like Delete but panic if an error occurs
func (k *ShardKind) MustDelete(ctx context.Context, key interface{}) *datastore.Key {
	_key, err := k.Delete(ctx, key)
	if err != nil {
		panic(err)
	}
	return _key
}

// DeleteMulti do Delete with multiple keys
func (k *ShardKind) DeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	var logger = xlog.WithContext(ctx).WithKey(ShardKindLoggerKey)
	var dsKeys, err = k.normMultiKeys(ctx, keys)
	if err != nil {
		return nil, err
	}
	var size = len(dsKeys)
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
		return nil, err
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
			"CounterShared#DeleteMulti [Datastore] (NoCache: %t)\n",
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
func (k *ShardKind) MustDeleteMulti(ctx context.Context, keys interface{}) []*datastore.Key {
	_keys, err := k.DeleteMulti(ctx, keys)
	if err != nil {
		panic(err)
	}
	return _keys
}

func (k *ShardKind) normMultiKeys(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	var dsKeys []*datastore.Key
	switch t := keys.(type) {
	case []string:
		tmp := keys.([]string)
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = helper.NewKey(ctx, "CounterShared", s)
		}
	case []interface{}:
		tmp := keys.([]interface{})
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = helper.NewKey(ctx, "CounterShared", s)
		}
	case []*datastore.Key:
		dsKeys = keys.([]*datastore.Key)
	default:
		return nil, fmt.Errorf("getmulti: unsupported keys type: %s", t)
	}
	return dsKeys, nil
}

// ShardQuery helps to build and execute a query
type ShardQuery struct {
	q *helper.Query
}

func NewShardQuery() *ShardQuery {
	return &ShardQuery{
		q: helper.NewQuery("CounterShared"),
	}
}

// Ancestor sets the ancestor filter
func (q *ShardQuery) Ancestor(a lazy.Value) *ShardQuery {
	q.q = q.q.Ancestor(a)
	return q
}

// Eq sets the "=" filter on the name field.
func (q *ShardQuery) Eq(name string, value lazy.Value) *ShardQuery {
	q.q = q.q.Eq(name, value)
	return q
}

// Lt sets the "<" filter on the "name" field.
func (q *ShardQuery) Lt(name string, value lazy.Value) *ShardQuery {
	q.q = q.q.Lt(name, value)
	return q
}

// Le sets the "<=" filter on the "name" field.
func (q *ShardQuery) Le(name string, value lazy.Value) *ShardQuery {
	q.q = q.q.Le(name, value)
	return q
}

// Gt sets the ">" filter on the "name" field.
func (q *ShardQuery) Gt(name string, value lazy.Value) *ShardQuery {
	q.q = q.q.Gt(name, value)
	return q
}

// Ge sets the ">=" filter on the "name" field.
func (q *ShardQuery) Ge(name string, value lazy.Value) *ShardQuery {
	q.q = q.q.Ge(name, value)
	return q
}

// Ne sets the "!=" filter on the "name" field.
func (q *ShardQuery) Ne(name string, value lazy.Value) *ShardQuery {
	q.q = q.q.Ne(name, value)
	return q
}

// Asc specifies ascending order on the given filed.
func (q *ShardQuery) Asc(name string) *ShardQuery {
	q.q = q.q.Asc(name)
	return q
}

// Desc specifies descending order on the given filed.
func (q *ShardQuery) Desc(name string) *ShardQuery {
	q.q = q.q.Desc(name)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ShardQuery) Limit(n lazy.Value) *ShardQuery {
	q.q = q.q.Limit(n)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ShardQuery) Start(value lazy.Value) *ShardQuery {
	q.q = q.q.Start(value)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ShardQuery) End(value lazy.Value) *ShardQuery {
	q.q = q.q.End(value)
	return q
}

// GetAll returns all key and value of the query.
func (q *ShardQuery) GetAll(ctx context.Context) ([]*datastore.Key, []*Shard, error) {
	var v []*Shard
	keys, err := q.q.GetAll(ctx, &v)
	if err != nil {
		return nil, nil, err
	}
	return keys, v, err
}

// MustGetAll is like GetAll but panic if an error occurrs.
func (q *ShardQuery) MustGetAll(ctx context.Context) ([]*datastore.Key, []*Shard) {
	keys, values, err := q.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	return keys, values
}

// GetAllValues is like GetAll but returns only values
func (q *ShardQuery) GetAllValues(ctx context.Context) ([]*Shard, error) {
	var v []*Shard
	_, err := q.q.GetAll(ctx, &v)
	if err != nil {
		return nil, err
	}
	return v, err
}

// MustGetAllValues is like GetAllValues but panic if an error occurrs
func (q *ShardQuery) MustGetAllValues(ctx context.Context) []*Shard {
	var v []*Shard
	_, err := q.q.GetAll(ctx, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// Count returns the count of entities
func (q *ShardQuery) Count(ctx context.Context) (int, error) {
	return q.q.Count(ctx)
}

// MustCount returns the count of entities
func (q *ShardQuery) MustCount(ctx context.Context) int {
	c, err := q.Count(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

type ShardPagination struct {
	Start string           `json:"start"`
	End   string           `json:"end"`
	Count int              `json:"count,omitempty"`
	Data  []*Shard         `json:"data"`
	Keys  []*datastore.Key `json:"-"`
}

// Run returns the a result as *ShardPagination object
func (q *ShardQuery) Run(ctx context.Context) (*ShardPagination, error) {
	iter, err := q.q.Run(ctx)
	if err != nil {
		return nil, err
	}
	pagination := &ShardPagination{}
	keys := []*datastore.Key{}
	data := []*Shard{}
	start, err := iter.Cursor()
	if err != nil {
		return nil, fmt.Errorf("couldn't get the start cursor: %v", err)
	}
	pagination.Start = start.String()
	for {
		var ent Shard
		key, err := iter.Next(&ent)
		if err == datastore.Done {
			end, err := iter.Cursor()
			if err != nil {
				return nil, fmt.Errorf("couldn't get the end cursor: %v", err)
			}
			pagination.Keys = keys
			pagination.Data = data
			pagination.End = end.String()
			return pagination, nil
		}
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		data = append(data, &ent)
	}
}

// MustRun is like Run but panic if an error occurrs
func (q *ShardQuery) MustRun(ctx context.Context) *ShardPagination {
	p, err := q.Run(ctx)
	if err != nil {
		panic(err)
	}
	return p
}
