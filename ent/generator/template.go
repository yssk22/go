package generator

const codeTemplate = `// Code generated by "ent -type={{.Type}}"; DO NOT EDIT

package {{.Package}}
{{ $v := .Instance }}
import(
    {{range $key, $as := .Dependencies -}}
    {{if $as -}}
    {{$as}} "{{$key}}"
    {{else -}}
    "{{$key}}"
    {{end -}}
    {{end }}
)

{{if .IsSearchable -}}
const {{.Type}}SearchIndexName = "ent.{{.Type}}"

// {{.Type}}SearchDoc is a object for search indexes.
type {{.Type}}SearchDoc struct {
	{{.IDField}} string // TODO: Support non-string key as ID
	{{range .Fields -}}{{if .IsSearch -}}
	{{.FieldName}} {{.SearchFieldTypeName}}
    {{end -}}{{end -}}
}

// ToSearchDoc returns a new *{{.Type}}SearchDoc
func ({{$v}} *{{.Type}}) ToSearchDoc() *{{.Type}}SearchDoc {
	s := &{{.Type}}SearchDoc{}
	s.{{.IDField}} = {{$v}}.{{.IDField}}
	{{range .Fields -}}{{if .IsSearch -}}
	{{ if .SearchFieldConverter -}}
	s.{{.FieldName}} = {{.SearchFieldConverter}}({{$v}}.{{.FieldName}})
	{{ else -}}
	s.{{.FieldName}} = {{$v}}.{{.FieldName}}
	{{end -}}{{end -}}{{end -}}
	return s
}

{{end -}}

func ({{$v}} *{{.Type}}) NewKey(ctx context.Context) *datastore.Key {
    return helper.NewKey(ctx, "{{.Kind}}", {{$v}}.{{.IDField}})
}

// UpdateByForm updates the fields by form values. All values should be validated
// before calling this function.
func ({{$v}} *{{.Type}}) UpdateByForm(form *keyvalue.GetProxy) {
    {{range .Fields -}}{{if .Form -}}
    if v, err := form.Get("{{.FieldNameSnakeCase}}"); err == nil {
        {{$v}}.{{.FieldName}} = {{.Form}}
    }
    {{end -}}{{end -}}
}


// New{{.Type}} returns a new *{{.Type}} with default field values.
func New{{.Type}}() *{{.Type}} {
    {{$v}} := &{{.Type}}{}
    {{range .Fields -}}{{if .Default -}}
    {{$v}}.{{.FieldName}} = {{.Default}}
    {{end}}{{end -}}
    return {{$v}}
}

type {{.Type}}Kind struct {
    BeforeSave func(ent *{{.Type}}) error
    AfterSave  func(ent *{{.Type}}) error
    useDefaultIfNil bool
    noCache bool
	noSearchIndexing bool
    noTimestampUpdate bool
}

// Default{{.Type}}Kind is a default value of *{{.Type}}Kind
var Default{{.Type}}Kind = &{{.Type}}Kind{}

// {{.Type}}KindLoggerKey is a logger key name for the ent
const {{.Type}}KindLoggerKey = "ent.{{snakecase .Kind}}"

func (k *{{.Type}}Kind) UseDefaultIfNil(b bool) *{{.Type}}Kind {
    k.useDefaultIfNil = b
    return k
}

// Get gets the kind entity from datastore
func (k *{{.Type}}Kind) Get(ctx context.Context, key interface{}) (*datastore.Key, *{{.Type}}, error) {
    keys, ents, err := k.GetMulti(ctx, []interface{}{key})
    if err != nil {
        return nil, nil, err
    }
    return keys[0], ents[0], nil
}

// MustGet is like Get but returns only values and panic if error happens.
func (k *{{.Type}}Kind) MustGet(ctx context.Context, key interface{}) *{{.Type}} {
    _, v, err := k.Get(ctx, key)
    if err != nil {
        panic(err)
    }
    return v
}

// GetMulti do Get with multiple keys. keys must be []string, []*datastore.Key, or []interface{}
func (k *{{.Type}}Kind) GetMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, []*{{.Type}}, error) {
    var logger = xlog.WithContext(ctx).WithKey({{.Type}}KindLoggerKey)
    var dsKeys, err = k.normMultiKeys(ctx, keys)
    if err != nil {
        return nil, nil, err
    }
    var size = len(dsKeys)
    var memKeys []string
    var ents []*{{.Type}}
    if size == 0 {
        return nil, nil, nil
    }
    ents = make([]*{{.Type}}, size, size)
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
        logger.Debug(func(p *xlog.Printer){
            p.Println("{{.Kind}}#GetMulti [Memcache]", )
            for i:=0; i < size; i++ {
                s := fmt.Sprintf("%v", ents[i])
                if len(s) > 20 {
                    p.Printf("\t%s - %s...\n", memKeys[i], s[:20])
                }else{
                    p.Printf("\t%s - %s\n", memKeys[i], s)
                }
                if i >= 20 {
                    p.Printf("\t...(and %d ents)\n", size - i)
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
    cacheMissingEnts := make([]*{{.Type}}, cacheMissingSize, cacheMissingSize)
    err = helper.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts)
    if helper.IsDatastoreError(err) {
        // we return nil even some ents hits the cache.
        return nil, nil, err
    }

    if k.useDefaultIfNil {
        for i:=0; i < cacheMissingSize; i++ {
            if cacheMissingEnts[i] == nil {
                cacheMissingEnts[i] = New{{.Type}}()
                cacheMissingEnts[i].{{.IDField}} = cacheMissingKeys[i].StringID() // TODO: Support non-string key as ID
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
        cacheEnts := make([]*{{.Type}}, 0)
        cacheKeys := make([]string, 0)
        for i := range ents {
            if ents[i] != nil {
                cacheEnts = append(cacheEnts, ents[i])
                cacheKeys = append(cacheKeys, memKeys[i])
            }
        }
        if len(cacheEnts) > 0 {
            if err := memcache.SetMulti(ctx, cacheKeys, cacheEnts); err != nil {
                logger.Warnf("Failed to create {{.Kind}}) caches: %v", err)
            }
        }
    }

    logger.Debug(func(p *xlog.Printer){
        p.Printf(
            "{{.Kind}}#GetMulti [Datastore] (UseDefault: %t, NoCache: %t)\n",
            k.useDefaultIfNil, k.noCache,
        )
        for i:=0; i < size; i++ {
            s := fmt.Sprintf("%v", ents[i])
            if len(s) > 20 {
                p.Printf("\t%s - %s...\n", dsKeys[i], s[:20])
            }else{
                p.Printf("\t%s - %s\n", dsKeys[i], s)
            }
            if i >= 20 {
                p.Printf("\t...(and %d ents)\n", size - i)
                break
            }
        }
    })
    return dsKeys, ents, nil
}

// MustGetMulti is like GetMulti but returns only values and panic if error happens.
func (k *{{.Type}}Kind) MustGetMulti(ctx context.Context, keys interface{}) []*{{.Type}} {
    _, v, err := k.GetMulti(ctx, keys)
    if err != nil {
        panic(err)
    }
    return v
}

// Put puts the entity to datastore.
func (k *{{.Type}}Kind) Put(ctx context.Context, ent *{{.Type}}) (*datastore.Key, error) {
    keys, err := k.PutMulti(ctx, []*{{.Type}}{
        ent,
    })
    if err != nil {
        return nil, err
    }
    return keys[0], nil
}

// MustPut is like Put and panic if an error occurrs.
func (k *{{.Type}}Kind) MustPut(ctx context.Context, ent *{{.Type}}) *datastore.Key {
    keys, err := k.Put(ctx, ent)
    if err != nil {
        panic(err)
    }
    return keys
}

// PutMulti do Put with multiple keys
func (k *{{.Type}}Kind) PutMulti(ctx context.Context, ents []*{{.Type}}) ([]*datastore.Key, error) {
    var size = len(ents)
    var dsKeys []*datastore.Key
	{{if .IsSearchable -}}
	var searchDocs []interface{} // to adopt search.Index#PutMulti()
	var searchKeys []string
	{{end -}}
    if size == 0 {
        return nil, nil
    }
	if size >= ent.MaxEntsPerPutDelete {
		return nil, ent.ErrTooManyEnts
	}
    logger := xlog.WithContext(ctx).WithKey({{.Type}}KindLoggerKey)

    dsKeys = make([]*datastore.Key, size, size)
    for i := range ents {
        if k.BeforeSave != nil {
            if err := k.BeforeSave(ents[i]); err != nil {
                return nil, err
            }
        }
        dsKeys[i] = ents[i].NewKey(ctx)
    }

	{{if .IsSearchable -}}
	if !k.noSearchIndexing {
		searchKeys := make([]string, size, size)
		searchDocs = make([]interface{}, size, size)
		for i := range ents {
			searchKeys[i] = dsKeys[i].Encode()
			searchDocs[i] = ents[i].ToSearchDoc()
		}
	}
	{{end -}}

    if !k.noTimestampUpdate {
        for i := range ents {
            ents[i].{{.TimestampField}} = xtime.Now()
        }
    }

    _, err := helper.PutMulti(ctx, dsKeys, ents)
    if helper.IsDatastoreError(err) {
        return nil, err
    }

    if !k.noCache {
        memKeys := make([]string, size, size)
        for i := range memKeys {
            memKeys[i] =ent.GetMemcacheKey(dsKeys[i])
        }
        err := memcache.DeleteMulti(ctx, memKeys)
        if memcache.IsMemcacheError(err) {
            logger.Warnf("Failed to invalidate memcache keys: %v", err)
        }
    }

	{{if .IsSearchable -}}
	if !k.noSearchIndexing {
		// TODO: should limit 200 docs per a call
		// see https://github.com/golang/appengine/blob/master/search/search.go#L136-L147
		index, err := search.Open({{.Type}}SearchIndexName)
		if err != nil {
            logger.Warnf("Failed to create search indexes (could not open index): %v ", err)
		} else {
			_, err = index.PutMulti(ctx, searchKeys, searchDocs)
			if err != nil {
	            logger.Warnf("Failed to create search indexes (PutMulti error): %v ", err)
			}
		}
	}
	{{end -}}

    logger.Debug(func(p *xlog.Printer){
        p.Printf(
            "{{.Kind}}#PutMulti [Datastore] (NoCache: %t)\n",
            k.noCache,
        )
        for i:=0; i < size; i++ {
            s := fmt.Sprintf("%v", ents[i])
            if len(s) > 20 {
                p.Printf("\t%s - %s...\n", dsKeys[i], s[:20])
            }else{
                p.Printf("\t%s - %s\n", dsKeys[i], s)
            }
            if i >= 20 {
                p.Printf("\t...(and %d ents)\n", size - i)
                break
            }
        }
    })

    return dsKeys, nil
}

// MustPutMulti is like PutMulti but panic if an error occurs
func (k *{{.Type}}Kind) MustPutMulti(ctx context.Context, ents []*{{.Type}}) ([]*datastore.Key) {
    keys, err := k.PutMulti(ctx, ents)
    if err != nil {
        panic(err)
    }
    return keys
}

// Delete deletes the entity from datastore
func (k *{{.Type}}Kind) Delete(ctx context.Context, key interface{}) (*datastore.Key, error) {
    keys, err := k.DeleteMulti(ctx, []interface{}{key})
    if err != nil {
        return nil, err
    }
    return keys[0], nil
}

// MustDelete is like Delete but panic if an error occurs
func (k *{{.Type}}Kind) MustDelete(ctx context.Context, key interface{}) (*datastore.Key) {
    _key, err := k.Delete(ctx, key)
    if err != nil {
        panic(err)
    }
    return _key
}

// DeleteMulti do Delete with multiple keys
func (k *{{.Type}}Kind) DeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
    var logger = xlog.WithContext(ctx).WithKey({{.Type}}KindLoggerKey)
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

	{{if .IsSearchable -}}
	var searchKeys []string
	if !k.noSearchIndexing {
		searchKeys = make([]string, size, size)
		for i, k := range dsKeys {
			searchKeys[i] = k.Encode()
		}
	}
	{{end -}}

    // Datastore access
    err = helper.DeleteMulti(ctx, dsKeys)
    if helper.IsDatastoreError(err) {
        // we return nil even some ents hits the cache.
        return nil, err
    }

    if !k.noCache {
        memKeys := make([]string, size, size)
        for i := range memKeys {
            memKeys[i] =ent.GetMemcacheKey(dsKeys[i])
        }
        err = memcache.DeleteMulti(ctx, memKeys)
        if memcache.IsMemcacheError(err) {
            logger.Warnf("Failed to invalidate memcache keys: %v", err)
        }
    }

	{{if .IsSearchable -}}
	if !k.noSearchIndexing {
		// TODO: should limit 200 docs per a call
		// see https://github.com/golang/appengine/blob/master/search/search.go#L136-L147
		index, err := search.Open({{.Type}}SearchIndexName)
		if err != nil {
            logger.Warnf("Failed to delete search indexes (could not open index): %v ", err)
		} else {
			err = index.DeleteMulti(ctx, searchKeys)
			if err != nil {
	            logger.Warnf("Failed to delete search indexes (PutMulti error): %v ", err)
			}
		}
	}
	{{end -}}


    logger.Debug(func(p *xlog.Printer){
        p.Printf(
            "{{.Kind}}#DeleteMulti [Datastore] (NoCache: %t)\n",
            k.noCache,
        )
        for i:=0; i < size; i++ {
            p.Printf("\t%s\n", dsKeys[i])
            if i >= 20 {
                p.Printf("\t...(and %d ents)\n", size - i)
                break
            }
        }
    })
    return dsKeys, nil
}

// MustDeleteMulti is like DeleteMulti but panic if an error occurs
func (k *{{.Type}}Kind) MustDeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key) {
    _keys, err := k.DeleteMulti(ctx, keys)
    if err != nil {
        panic(err)
    }
    return _keys
}

func (k *{{.Type}}Kind) normMultiKeys(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
    var dsKeys []*datastore.Key
    switch t := keys.(type) {
        case []string:
            tmp := keys.([]string)
            dsKeys = make([]*datastore.Key, len(tmp))
            for i, s := range tmp {
                dsKeys[i] = helper.NewKey(ctx, "{{.Kind}}", s)
            }
        case []interface{}:
            tmp := keys.([]interface{})
            dsKeys = make([]*datastore.Key, len(tmp))
            for i, s := range tmp {
                dsKeys[i] = helper.NewKey(ctx, "{{.Kind}}", s)
            }
        case []*datastore.Key:
            dsKeys = keys.([]*datastore.Key)
        default:
            return nil, fmt.Errorf("getmulti: unsupported keys type: %s", t)
    }
    return dsKeys, nil
}

// {{.Type}}Query helps to build and execute a query
type {{.Type}}Query struct {
    q *helper.Query
}

func New{{.Type}}Query() *{{.Type}}Query {
    return &{{.Type}}Query{
        q: helper.NewQuery("{{.Kind}}"),
    }
}

// Ancestor sets the ancestor filter
func (q *{{.Type}}Query) Ancestor(a lazy.Value) *{{.Type}}Query {
	q.q = q.q.Ancestor(a)
	return q
}

// Eq sets the "=" filter on the name field.
func (q *{{.Type}}Query) Eq(name string, value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Eq(name, value)
	return q
}

// Lt sets the "<" filter on the "name" field.
func (q *{{.Type}}Query) Lt(name string, value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Lt(name, value)
	return q
}

// Le sets the "<=" filter on the "name" field.
func (q *{{.Type}}Query) Le(name string, value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Le(name, value)
	return q
}

// Gt sets the ">" filter on the "name" field.
func (q *{{.Type}}Query) Gt(name string, value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Gt(name, value)
	return q
}

// Ge sets the ">=" filter on the "name" field.
func (q *{{.Type}}Query) Ge(name string, value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Ge(name, value)
	return q
}

// Ne sets the "!=" filter on the "name" field.
func (q *{{.Type}}Query) Ne(name string, value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Ne(name, value)
	return q
}

// Asc specifies ascending order on the given filed.
func (q *{{.Type}}Query) Asc(name string) *{{.Type}}Query {
	q.q = q.q.Asc(name)
	return q
}

// Desc specifies descending order on the given filed.
func (q *{{.Type}}Query) Desc(name string) *{{.Type}}Query {
	q.q = q.q.Desc(name)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *{{.Type}}Query) Limit(n lazy.Value) *{{.Type}}Query {
	q.q = q.q.Limit(n)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *{{.Type}}Query) Start(value lazy.Value) *{{.Type}}Query {
	q.q = q.q.Start(value)
	return q
}

// Limit specifies the numbe of limit returend by this query.
func (q *{{.Type}}Query) End(value lazy.Value) *{{.Type}}Query {
	q.q = q.q.End(value)
	return q
}

// GetAll returns all key and value of the query.
func (q *{{.Type}}Query) GetAll(ctx context.Context) ([]*datastore.Key, []*{{.Type}}, error) {
    var v []*{{.Type}}
    keys, err := q.q.GetAll(ctx, &v)
    if err != nil {
        return nil, nil, err
    }
    return keys, v, err
}

// MustGetAll is like GetAll but panic if an error occurrs.
func (q *{{.Type}}Query) MustGetAll(ctx context.Context) ([]*datastore.Key, []*{{.Type}}) {
    keys, values, err := q.GetAll(ctx)
    if err != nil {
        panic(err)
    }
    return keys, values
}

// GetAllValues is like GetAll but returns only values
func (q *{{.Type}}Query) GetAllValues(ctx context.Context) ([]*{{.Type}}, error) {
    var v []*{{.Type}}
    _, err := q.q.GetAll(ctx, &v)
    if err != nil {
        return nil, err
    }
    return v, err
}

// MustGetAllValues is like GetAllValues but panic if an error occurrs
func (q *{{.Type}}Query) MustGetAllValues(ctx context.Context) []*{{.Type}} {
    var v []*{{.Type}}
    _, err := q.q.GetAll(ctx, &v)
    if err != nil {
        panic(err)
    }
    return v
}

// Count returns the count of entities
func (q *{{.Type}}Query) Count(ctx context.Context) (int, error) {
    return q.q.Count(ctx)
}

// MustCount returns the count of entities
func (q *{{.Type}}Query) MustCount(ctx context.Context) int {
    c, err := q.Count(ctx)
    if err != nil {
        panic(err)
    }
    return c
}

type {{.Type}}Pagination struct {` +
	"Start string           `json:\"start\"`\n" +
	"End   string           `json:\"end\"`\n" +
	"Count int              `json:\"count,omitempty\"`\n" +
	"Data  []*{{.Type}}     `json:\"data\"`\n" +
	"Keys  []*datastore.Key `json:\"-\"`\n" + `
}

// Run returns the a result as *{{.Type}}Pagination object
func (q *{{.Type}}Query) Run(ctx context.Context) (*{{.Type}}Pagination, error) {
    iter, err := q.q.Run(ctx)
    if err != nil {
        return nil, err
    }
    pagination := &{{.Type}}Pagination{}
    keys := []*datastore.Key{}
    data := []*{{.Type}}{}
	start, err := iter.Cursor()
	if err != nil {
		return nil, fmt.Errorf("couldn't get the start cursor: %v", err)
	}
	pagination.Start = start.String()
    for {
        var ent {{.Type}}
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
func (q *{{.Type}}Query) MustRun(ctx context.Context) *{{.Type}}Pagination {
    p, err := q.Run(ctx)
    if err != nil {
        panic(err)
    }
    return p
}
`
