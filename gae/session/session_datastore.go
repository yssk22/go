// Code generated by "ent -type=Session"; DO NOT EDIT

package session

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

func (s *Session) NewKey(ctx context.Context) *datastore.Key {
	return helper.NewKey(ctx, "Session", s.ID)
}

// UpdateByForm updates the fields by form values. All values should be validated
// before calling this function.
func (s *Session) UpdateByForm(form *keyvalue.GetProxy) {
}

// NewSession returns a new *Session with default field values.
func NewSession() *Session {
	s := &Session{}
	return s
}

type SessionKind struct {
	BeforeSave        func(ent *Session) error
	AfterSave         func(ent *Session) error
	useDefaultIfNil   bool
	noCache           bool
	noTimestampUpdate bool
}

// DefaultSessionKind is a default value of *SessionKind
var DefaultSessionKind = &SessionKind{}

const SessionKindLoggerKey = "ent.session"

func (k *SessionKind) UseDefaultIfNil(b bool) *SessionKind {
	k.useDefaultIfNil = b
	return k
}

// Get gets the kind entity from datastore
func (k *SessionKind) Get(ctx context.Context, key interface{}) (*datastore.Key, *Session, error) {
	keys, ents, err := k.GetMulti(ctx, key)
	if err != nil {
		return nil, nil, err
	}
	return keys[0], ents[0], nil
}

// MustGet is like Get but returns only values and panic if error happens.
func (k *SessionKind) MustGet(ctx context.Context, key interface{}) *Session {
	_, v, err := k.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMulti do Get with multiple keys
func (k *SessionKind) GetMulti(ctx context.Context, keys ...interface{}) ([]*datastore.Key, []*Session, error) {
	var size = len(keys)
	var memKeys []string
	var dsKeys []*datastore.Key
	var ents []*Session
	if size == 0 {
		return nil, nil, nil
	}
	logger := xlog.WithContext(ctx).WithKey(SessionKindLoggerKey)
	dsKeys = make([]*datastore.Key, size, size)
	for i := range keys {
		dsKeys[i] = helper.NewKey(ctx, "Session", keys[i])
	}
	ents = make([]*Session, size, size)
	// Memcache access
	if !k.noCache {
		logger.Debugf("Trying to get entities from memcache...")
		memKeys = make([]string, size, size)
		for i := range dsKeys {
			memKeys[i] = ent.GetMemcacheKey(dsKeys[i])
		}
		err := memcache.GetMulti(ctx, memKeys, ents)
		if err == nil {
			// Hit caches on all keys!!
			return dsKeys, ents, nil
		}
		logger.Debug(func(p *xlog.Printer) {
			p.Println("Session#GetMulti [Memcache]")
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
	cacheMissingEnts := make([]*Session, cacheMissingSize, cacheMissingSize)
	err := helper.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts)
	if helper.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, nil, err
	}

	if k.useDefaultIfNil {
		for i := 0; i < cacheMissingSize; i++ {
			if cacheMissingEnts[i] == nil {
				cacheMissingEnts[i] = NewSession()
				cacheMissingEnts[i].ID = dsKeys[i].StringID() // TODO: Support non-string key as ID
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
		cacheEnts := make([]*Session, 0)
		cacheKeys := make([]string, 0)
		for i := range ents {
			if ents[i] != nil {
				cacheEnts = append(cacheEnts, ents[i])
				cacheKeys = append(cacheKeys, memKeys[i])
			}
		}
		if len(cacheEnts) > 0 {
			if err := memcache.SetMulti(ctx, cacheKeys, cacheEnts); err != nil {
				logger.Warnf("Failed to create Session) caches: %v", err)
			}
		}
	}

	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"Session#GetMulti [Datastore] (UseDefault: %t, NoCache: %t)\n",
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
func (k *SessionKind) MustGetMulti(ctx context.Context, keys ...interface{}) []*Session {
	_, v, err := k.GetMulti(ctx, keys...)
	if err != nil {
		panic(err)
	}
	return v
}

// Put puts the entity to datastore.
func (k *SessionKind) Put(ctx context.Context, ent *Session) (*datastore.Key, error) {
	keys, err := k.PutMulti(ctx, []*Session{
		ent,
	})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustPut is like Put and panic if an error occurrs.
func (k *SessionKind) MustPut(ctx context.Context, ent *Session) *datastore.Key {
	keys, err := k.Put(ctx, ent)
	if err != nil {
		panic(err)
	}
	return keys
}

// PutMulti do Put with multiple keys
func (k *SessionKind) PutMulti(ctx context.Context, ents []*Session) ([]*datastore.Key, error) {
	var size = len(ents)
	var dsKeys []*datastore.Key
	if size == 0 {
		return nil, nil
	}
	logger := xlog.WithContext(ctx).WithKey(SessionKindLoggerKey)

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
			ents[i].Timestamp = xtime.Now()
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
			"Session#PutMulti [Datastore] (NoCache: %t)\n",
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
func (k *SessionKind) MustPutMulti(ctx context.Context, ents []*Session) []*datastore.Key {
	keys, err := k.PutMulti(ctx, ents)
	if err != nil {
		panic(err)
	}
	return keys
}

// Delete deletes the entity from datastore
func (k *SessionKind) Delete(ctx context.Context, key interface{}) (*datastore.Key, error) {
	keys, err := k.DeleteMulti(ctx, key)
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustDelete is like Delete but panic if an error occurs
func (k *SessionKind) MustDelete(ctx context.Context, key interface{}) *datastore.Key {
	keys, err := k.DeleteMulti(ctx, key)
	if err != nil {
		panic(err)
	}
	return keys[0]
}

// DeleteMulti do Delete with multiple keys
func (k *SessionKind) DeleteMulti(ctx context.Context, keys ...interface{}) ([]*datastore.Key, error) {
	var size = len(keys)
	var dsKeys []*datastore.Key
	if size == 0 {
		return nil, nil
	}
	logger := xlog.WithContext(ctx).WithKey(SessionKindLoggerKey)
	dsKeys = make([]*datastore.Key, size, size)
	for i := range keys {
		dsKeys[i] = helper.NewKey(ctx, "Session", keys[i])
	}
	// Datastore access
	err := helper.DeleteMulti(ctx, dsKeys)
	if helper.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
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
			"Session#DeleteMulti [Datastore] (NoCache: %t)\n",
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
func (k *SessionKind) MustDeleteMulti(ctx context.Context, keys ...interface{}) []*datastore.Key {
	_keys, err := k.DeleteMulti(ctx, keys...)
	if err != nil {
		panic(err)
	}
	return _keys
}

// SessionQuery helps to build and execute a query
type SessionQuery struct {
	q *helper.Query
}

func NewSessionQuery() *SessionQuery {
	return &SessionQuery{
		q: helper.NewQuery("Session"),
	}
}

// Ancestor sets the ancestor filter
func (q *SessionQuery) Ancestor(a lazy.Value) *SessionQuery {
	q.q = q.q.Ancestor(a)
	return q
}

// Eq sets the "=" filter on the name field.
func (q *SessionQuery) Eq(name string, value lazy.Value) *SessionQuery {
	q.q = q.q.Eq(name, value)
	return q
}

// Lt sets the "<" filter on the "name" field.
func (q *SessionQuery) Lt(name string, value lazy.Value) *SessionQuery {
	q.q = q.q.Lt(name, value)
	return q
}

// Le sets the "<=" filter on the "name" field.
func (q *SessionQuery) Le(name string, value lazy.Value) *SessionQuery {
	q.q = q.q.Le(name, value)
	return q
}

// Gt sets the ">" filter on the "name" field.
func (q *SessionQuery) Gt(name string, value lazy.Value) *SessionQuery {
	q.q = q.q.Gt(name, value)
	return q
}

// Ge sets the ">=" filter on the "name" field.
func (q *SessionQuery) Ge(name string, value lazy.Value) *SessionQuery {
	q.q = q.q.Ge(name, value)
	return q
}

// Ne sets the "!=" filter on the "name" field.
func (q *SessionQuery) Ne(name string, value lazy.Value) *SessionQuery {
	q.q = q.q.Ne(name, value)
	return q
}

// Asc specifies ascending order on the given filed.
func (q *SessionQuery) Asc(name string) *SessionQuery {
	q.q = q.q.Asc(name)
	return q
}

// Desc specifies descending order on the given filed.
func (q *SessionQuery) Desc(name string) *SessionQuery {
	q.q = q.q.Desc(name)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *SessionQuery) Limit(n lazy.Value) *SessionQuery {
	q.q = q.q.Limit(n)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *SessionQuery) Start(value lazy.Value) *SessionQuery {
	q.q = q.q.Start(value)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *SessionQuery) End(value lazy.Value) *SessionQuery {
	q.q = q.q.End(value)
	return q
}

// GetAll returns all key and value of the query.
func (q *SessionQuery) GetAll(ctx context.Context) ([]*datastore.Key, []*Session, error) {
	var v []*Session
	keys, err := q.q.GetAll(ctx, &v)
	if err != nil {
		return nil, nil, err
	}
	return keys, v, err
}

// MustGetAll is like GetAll but panic if an error occurrs.
func (q *SessionQuery) MustGetAll(ctx context.Context) ([]*datastore.Key, []*Session) {
	keys, values, err := q.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	return keys, values
}

// GetAllValues is like GetAll but returns only values
func (q *SessionQuery) GetAllValues(ctx context.Context) ([]*Session, error) {
	var v []*Session
	_, err := q.q.GetAll(ctx, &v)
	if err != nil {
		return nil, err
	}
	return v, err
}

// MustGetAllValues is like GetAllValues but panic if an error occurrs
func (q *SessionQuery) MustGetAllValues(ctx context.Context) []*Session {
	var v []*Session
	_, err := q.q.GetAll(ctx, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// Count returns the count of entities
func (q *SessionQuery) Count(ctx context.Context) (int, error) {
	return q.q.Count(ctx)
}

// MustCount returns the count of entities
func (q *SessionQuery) MustCount(ctx context.Context) int {
	c, err := q.Count(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

type SessionPagination struct {
	Start datastore.Cursor `json:"start"`
	End   datastore.Cursor `json:"end"`
	Data  []*Session       `json:"data"`
	Keys  []*datastore.Key `json:"-"`
}

// Run returns the a result as *SessionPagination object
func (q *SessionQuery) Run(ctx context.Context) (*SessionPagination, error) {
	iter, err := q.q.Run(ctx)
	if err != nil {
		return nil, err
	}
	start, err := iter.Cursor()
	if err != nil {
		return nil, fmt.Errorf("couldn't get the start cursor: %v", err)
	}
	var keys []*datastore.Key
	var data []*Session
	for {
		var ent Session
		key, err := iter.Next(&ent)
		if err == datastore.Done {
			end, err := iter.Cursor()
			if err != nil {
				return nil, fmt.Errorf("couldn't get the end cursor: %v", err)
			}
			return &SessionPagination{
				Start: start,
				End:   end,
				Keys:  keys,
				Data:  data,
			}, nil
		}
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		data = append(data, &ent)
	}
}

// MustRun is like Run but panic if an error occurrs
func (q *SessionQuery) MustRun(ctx context.Context) *SessionPagination {
	p, err := q.Run(ctx)
	if err != nil {
		panic(err)
	}
	return p
}
