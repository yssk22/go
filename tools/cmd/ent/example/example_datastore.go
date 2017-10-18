// Code generated by "ent -type=Example"; DO NOT EDIT

package example

import (
	"context"
	"fmt"
	"github.com/speedland/go/ent"
	helper "github.com/speedland/go/gae/datastore"
	"github.com/speedland/go/gae/memcache"
	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/rgb"
	"github.com/speedland/go/x/xerrors"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
)

const ExampleSearchIndexName = "ent.Example"

// ExampleSearchDoc is a object for search indexes.
type ExampleSearchDoc struct {
	ID           string // TODO: Support non-string key as ID
	Desc         string
	ContentBytes search.HTML
	BoolType     search.Atom
	FloatType    float64
	Location     appengine.GeoPoint
}

// ToSearchDoc returns a new *ExampleSearchDoc
func (e *Example) ToSearchDoc() *ExampleSearchDoc {
	s := &ExampleSearchDoc{}
	s.ID = e.ID
	s.Desc = e.Desc
	s.ContentBytes = ent.BytesToHTML(e.ContentBytes)
	s.BoolType = ent.BoolToAtom(e.BoolType)
	s.FloatType = e.FloatType
	s.Location = e.Location
	return s
}

func (e *Example) NewKey(ctx context.Context) *datastore.Key {
	return helper.NewKey(ctx, "Example", e.ID)
}

// UpdateByForm updates the fields by form values. All values should be validated
// before calling this function.
func (e *Example) UpdateByForm(form *keyvalue.GetProxy) {
	if v, err := form.Get("digit"); err == nil {
		e.Digit = ent.ParseInt(v.(string))
	}
	if v, err := form.Get("desc"); err == nil {
		e.Desc = v.(string)
	}
	if v, err := form.Get("content_bytes"); err == nil {
		e.ContentBytes = []byte(v.(string))
	}
	if v, err := form.Get("slice_type"); err == nil {
		e.SliceType = ent.ParseStringList(v.(string))
	}
	if v, err := form.Get("bool_type"); err == nil {
		e.BoolType = ent.ParseBool(v.(string))
	}
	if v, err := form.Get("float_type"); err == nil {
		e.FloatType = ent.ParseFloat64(v.(string))
	}
	if v, err := form.Get("custom_type"); err == nil {
		e.CustomType = rgb.MustParseRGB(v.(string))
	}
}

// NewExample returns a new *Example with default field values.
func NewExample() *Example {
	e := &Example{}
	e.Digit = 10
	e.Desc = "This is default value"
	e.CreatedAt = ent.ParseTime("$now")
	e.DefaultTime = ent.ParseTime("2016-01-01T20:12:10Z")
	return e
}

type ExampleKind struct {
	useDefaultIfNil           bool
	noCache                   bool
	noSearchIndexing          bool
	ignoreSearchIndexingError bool
	noTimestampUpdate         bool
	enforceNamespace          bool
	namespace                 string
}

type ExampleKindReplacer interface {
	Replace(*Example, *Example) *Example
}

type ExampleKindReplacerFunc func(*Example, *Example) *Example

func (f ExampleKindReplacerFunc) Replace(ent1 *Example, ent2 *Example) *Example {
	return f(ent1, ent2)
}

// DefaultExampleKind is a default value of *ExampleKind
var DefaultExampleKind = &ExampleKind{}

// ExampleKindLoggerKey is a logger key name for the ent
const ExampleKindLoggerKey = "ent.example"

// EnforceNamespace enforces namespace for Get/Put/Delete or not.
func (k *ExampleKind) EnforceNamespace(ns string, b bool) *ExampleKind {
	k.enforceNamespace = b
	k.namespace = ns
	return k
}

func (k *ExampleKind) UseDefaultIfNil(b bool) *ExampleKind {
	k.useDefaultIfNil = b
	return k
}

// Get gets the kind entity from datastore
func (k *ExampleKind) Get(ctx context.Context, key interface{}) (*datastore.Key, *Example, error) {
	keys, ents, err := k.GetMulti(ctx, []interface{}{key})
	if err != nil {
		return nil, nil, err
	}
	return keys[0], ents[0], nil
}

// MustGet is like Get but returns only values and panic if error happens.
func (k *ExampleKind) MustGet(ctx context.Context, key interface{}) *Example {
	_, v, err := k.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMulti do Get with multiple keys. keys must be []string, []*datastore.Key, or []interface{}
func (k *ExampleKind) GetMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, []*Example, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ExampleKindLoggerKey)
	var err error
	var dsKeys []*datastore.Key
	var memKeys []string
	var ents []*Example
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
	ents = make([]*Example, size, size)
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
			p.Println("Example#GetMulti [Memcache]")
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
	cacheMissingEnts := make([]*Example, cacheMissingSize, cacheMissingSize)
	err = helper.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts)
	if helper.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, nil, xerrors.Wrap(err, "datastore error")
	}

	if k.useDefaultIfNil {
		for i := 0; i < cacheMissingSize; i++ {
			if cacheMissingEnts[i] == nil {
				cacheMissingEnts[i] = NewExample()
				cacheMissingEnts[i].ID = cacheMissingKeys[i].StringID() // TODO: Support non-string key as ID
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
		cacheEnts := make([]*Example, 0)
		cacheKeys := make([]string, 0)
		for i := range ents {
			if ents[i] != nil {
				cacheEnts = append(cacheEnts, ents[i])
				cacheKeys = append(cacheKeys, memKeys[i])
			}
		}
		if len(cacheEnts) > 0 {
			if err := memcache.SetMulti(ctx, cacheKeys, cacheEnts); err != nil {
				logger.Warnf("Failed to create Example) caches: %v", err)
			}
		}
	}

	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"Example#GetMulti [Datastore] (UseDefault: %t, NoCache: %t)\n",
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
func (k *ExampleKind) MustGetMulti(ctx context.Context, keys interface{}) []*Example {
	_, v, err := k.GetMulti(ctx, keys)
	if err != nil {
		panic(err)
	}
	return v
}

// Put puts the entity to datastore.
func (k *ExampleKind) Put(ctx context.Context, ent *Example) (*datastore.Key, error) {
	keys, err := k.PutMulti(ctx, []*Example{
		ent,
	})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustPut is like Put and panic if an error occurrs.
func (k *ExampleKind) MustPut(ctx context.Context, ent *Example) *datastore.Key {
	keys, err := k.Put(ctx, ent)
	if err != nil {
		panic(err)
	}
	return keys
}

// PutMulti do Put with multiple keys
func (k *ExampleKind) PutMulti(ctx context.Context, ents []*Example) ([]*datastore.Key, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ExampleKindLoggerKey)
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	var searchDocs []interface{} // to adopt search.Index#PutMulti()
	var searchKeys []string
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

	if !k.noSearchIndexing {
		searchKeys = make([]string, size, size)
		searchDocs = make([]interface{}, size, size)
		for i := range ents {
			searchKeys[i] = dsKeys[i].Encode()
			searchDocs[i] = ents[i].ToSearchDoc()
		}
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

	if !k.noSearchIndexing {
		// TODO: should limit 200 docs per a call
		// see https://github.com/golang/appengine/blob/master/search/search.go#L136-L147
		index, err := search.Open(ExampleSearchIndexName)
		if err != nil {
			err = xerrors.Wrap(err, "search.Open(%q) returns errors", ExampleSearchIndexName)
			if !k.ignoreSearchIndexingError {
				return nil, err
			} else {
				logger.Warnf(err.Error())
			}
		} else {
			_, err = index.PutMulti(ctx, searchKeys, searchDocs)
			if err != nil {
				err = xerrors.Wrap(err, "index.PutMulti returns errors")
				if !k.ignoreSearchIndexingError {
					return nil, err
				} else {
					logger.Warnf(err.Error())
				}
			}
		}
	}
	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"Example#PutMulti [Datastore] (NoCache: %t)\n",
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
func (k *ExampleKind) MustPutMulti(ctx context.Context, ents []*Example) []*datastore.Key {
	keys, err := k.PutMulti(ctx, ents)
	if err != nil {
		panic(err)
	}
	return keys
}

func (k *ExampleKind) Replace(ctx context.Context, ent *Example, replacer ExampleKindReplacer) (*datastore.Key, *Example, error) {
	keys, ents, err := k.ReplaceMulti(ctx, []*Example{
		ent,
	}, replacer)
	if err != nil {
		return nil, ents[0], err
	}
	return keys[0], ents[0], err
}

func (k *ExampleKind) MustReplace(ctx context.Context, ent *Example, replacer ExampleKindReplacer) (*datastore.Key, *Example) {
	key, ent, err := k.Replace(ctx, ent, replacer)
	if err != nil {
		panic(err)
	}
	return key, ent
}

func (k *ExampleKind) ReplaceMulti(ctx context.Context, ents []*Example, replacer ExampleKindReplacer) ([]*datastore.Key, []*Example, error) {
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

func (k *ExampleKind) MustReplaceMulti(ctx context.Context, ents []*Example, replacer ExampleKindReplacer) ([]*datastore.Key, []*Example) {
	keys, ents, err := k.ReplaceMulti(ctx, ents, replacer)
	if err != nil {
		panic(err)
	}
	return keys, ents
}

// Delete deletes the entity from datastore
func (k *ExampleKind) Delete(ctx context.Context, key interface{}) (*datastore.Key, error) {
	keys, err := k.DeleteMulti(ctx, []interface{}{key})
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

// MustDelete is like Delete but panic if an error occurs
func (k *ExampleKind) MustDelete(ctx context.Context, key interface{}) *datastore.Key {
	_key, err := k.Delete(ctx, key)
	if err != nil {
		panic(err)
	}
	return _key
}

// DeleteMulti do Delete with multiple keys
func (k *ExampleKind) DeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ExampleKindLoggerKey)
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

	var searchKeys []string
	if !k.noSearchIndexing {
		searchKeys = make([]string, size, size)
		for i, k := range dsKeys {
			searchKeys[i] = k.Encode()
		}
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

	if !k.noSearchIndexing {
		// TODO: should limit 200 docs per a call
		// see https://github.com/golang/appengine/blob/master/search/search.go#L136-L147
		index, err := search.Open(ExampleSearchIndexName)
		if err != nil {
			logger.Warnf("Failed to delete search indexes (could not open index): %v ", err)
		} else {
			err = index.DeleteMulti(ctx, searchKeys)
			if err != nil {
				logger.Warnf("Failed to delete search indexes (PutMulti error): %v ", err)
			}
		}
	}
	logger.Debug(func(p *xlog.Printer) {
		p.Printf(
			"Example#DeleteMulti [Datastore] (NoCache: %t)\n",
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
func (k *ExampleKind) MustDeleteMulti(ctx context.Context, keys interface{}) []*datastore.Key {
	_keys, err := k.DeleteMulti(ctx, keys)
	if err != nil {
		panic(err)
	}
	return _keys
}

func (k *ExampleKind) normMultiKeys(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	var dsKeys []*datastore.Key
	switch t := keys.(type) {
	case []string:
		tmp := keys.([]string)
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = helper.NewKey(ctx, "Example", s)
		}
	case []interface{}:
		tmp := keys.([]interface{})
		dsKeys = make([]*datastore.Key, len(tmp))
		for i, s := range tmp {
			dsKeys[i] = helper.NewKey(ctx, "Example", s)
		}
	case []*datastore.Key:
		dsKeys = keys.([]*datastore.Key)
	default:
		return nil, fmt.Errorf("unsupported keys type: %s", t)
	}
	return dsKeys, nil
}

// ExampleQuery helps to build and execute a query
type ExampleQuery struct {
	q *helper.Query
}

func NewExampleQuery() *ExampleQuery {
	return &ExampleQuery{
		q: helper.NewQuery("Example"),
	}
}

// Ancestor sets the ancestor filter
func (q *ExampleQuery) Ancestor(a lazy.Value) *ExampleQuery {
	q.q = q.q.Ancestor(a)
	return q
}

// Eq sets the "=" filter on the name field.
func (q *ExampleQuery) Eq(name string, value lazy.Value) *ExampleQuery {
	q.q = q.q.Eq(name, value)
	return q
}

// Lt sets the "<" filter on the "name" field.
func (q *ExampleQuery) Lt(name string, value lazy.Value) *ExampleQuery {
	q.q = q.q.Lt(name, value)
	return q
}

// Le sets the "<=" filter on the "name" field.
func (q *ExampleQuery) Le(name string, value lazy.Value) *ExampleQuery {
	q.q = q.q.Le(name, value)
	return q
}

// Gt sets the ">" filter on the "name" field.
func (q *ExampleQuery) Gt(name string, value lazy.Value) *ExampleQuery {
	q.q = q.q.Gt(name, value)
	return q
}

// Ge sets the ">=" filter on the "name" field.
func (q *ExampleQuery) Ge(name string, value lazy.Value) *ExampleQuery {
	q.q = q.q.Ge(name, value)
	return q
}

// Ne sets the "!=" filter on the "name" field.
func (q *ExampleQuery) Ne(name string, value lazy.Value) *ExampleQuery {
	q.q = q.q.Ne(name, value)
	return q
}

// Asc specifies ascending order on the given filed.
func (q *ExampleQuery) Asc(name string) *ExampleQuery {
	q.q = q.q.Asc(name)
	return q
}

// Desc specifies descending order on the given filed.
func (q *ExampleQuery) Desc(name string) *ExampleQuery {
	q.q = q.q.Desc(name)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ExampleQuery) Limit(n lazy.Value) *ExampleQuery {
	q.q = q.q.Limit(n)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ExampleQuery) Start(value lazy.Value) *ExampleQuery {
	q.q = q.q.Start(value)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *ExampleQuery) End(value lazy.Value) *ExampleQuery {
	q.q = q.q.End(value)
	return q
}

// GetAll returns all key and value of the query.
func (q *ExampleQuery) GetAll(ctx context.Context) ([]*datastore.Key, []*Example, error) {
	var v []*Example
	keys, err := q.q.GetAll(ctx, &v)
	if err != nil {
		return nil, nil, err
	}
	return keys, v, err
}

// MustGetAll is like GetAll but panic if an error occurrs.
func (q *ExampleQuery) MustGetAll(ctx context.Context) ([]*datastore.Key, []*Example) {
	keys, values, err := q.GetAll(ctx)
	if err != nil {
		panic(err)
	}
	return keys, values
}

// GetAllValues is like GetAll but returns only values
func (q *ExampleQuery) GetAllValues(ctx context.Context) ([]*Example, error) {
	var v []*Example
	_, err := q.q.GetAll(ctx, &v)
	if err != nil {
		return nil, err
	}
	return v, err
}

// MustGetAllValues is like GetAllValues but panic if an error occurrs
func (q *ExampleQuery) MustGetAllValues(ctx context.Context) []*Example {
	var v []*Example
	_, err := q.q.GetAll(ctx, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// Count returns the count of entities
func (q *ExampleQuery) Count(ctx context.Context) (int, error) {
	return q.q.Count(ctx)
}

// MustCount returns the count of entities
func (q *ExampleQuery) MustCount(ctx context.Context) int {
	c, err := q.Count(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

type ExamplePagination struct {
	Start string           `json:"start"`
	End   string           `json:"end"`
	Count int              `json:"count,omitempty"`
	Data  []*Example       `json:"data"`
	Keys  []*datastore.Key `json:"-"`
}

// Run returns the a result as *ExamplePagination object
func (q *ExampleQuery) Run(ctx context.Context) (*ExamplePagination, error) {
	iter, err := q.q.Run(ctx)
	if err != nil {
		return nil, err
	}
	pagination := &ExamplePagination{}
	keys := []*datastore.Key{}
	data := []*Example{}
	start, err := iter.Cursor()
	if err != nil {
		return nil, fmt.Errorf("couldn't get the start cursor: %v", err)
	}
	pagination.Start = start.String()
	for {
		var ent Example
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
func (q *ExampleQuery) MustRun(ctx context.Context) *ExamplePagination {
	p, err := q.Run(ctx)
	if err != nil {
		panic(err)
	}
	return p
}

// SearchKeys returns the a result as *ExamplePagination object. It only containd valid []*datastore.Keys
func (k *ExampleKind) SearchKeys(ctx context.Context, query string, opts *search.SearchOptions) (*ExamplePagination, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ExampleKindLoggerKey)
	index, err := search.Open(ExampleSearchIndexName)
	if err != nil {
		return nil, err
	}
	// we don't need to grab document data since we can grab documents from datastore.
	if opts == nil {
		opts = &search.SearchOptions{}
	}
	opts.IDsOnly = true
	iter := index.Search(ctx, query, opts)
	pagination := &ExamplePagination{}
	keys := []*datastore.Key{}
	pagination.Start = string(iter.Cursor())
	pagination.Count = iter.Count()
	for {
		var ent ExampleSearchDoc
		id, err := iter.Next(&ent)
		if err == search.Done {
			pagination.Keys = keys
			pagination.End = string(iter.Cursor())
			return pagination, nil
		}
		if err != nil {
			return nil, err
		}
		key, err := datastore.DecodeKey(id)
		if err != nil {
			logger.Warnf("unexpected search doc id found: %s (%s), skipping...", id, err)
		} else {
			keys = append(keys, key)
		}
	}
}

// SearchValues returns the a result as *ExamplePagination object with filling Data field.
func (k *ExampleKind) SearchValues(ctx context.Context, query string, opts *search.SearchOptions) (*ExamplePagination, error) {
	ctx, logger := xlog.WithContextAndKey(ctx, "", ExampleKindLoggerKey)
	p, err := k.SearchKeys(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	keys, values, err := k.GetMulti(ctx, p.Keys)
	if err != nil {
		return nil, err
	}
	for i, v := range values {
		if v == nil {
			logger.Warnf("search index found, but no ent found - (Kind:Example, Key:%s)", keys[i].StringID())
		} else {
			p.Data = append(p.Data, v)
		}
	}
	return p, nil
}

// DeleteMatched deletes the all ents that match with the query.
// This func modify Limit/StartKey condition in the query so that you should restore it
// if you want to reuse the query.
func (k *ExampleKind) DeleteMatched(ctx context.Context, q *ExampleQuery) (int, error) {
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
