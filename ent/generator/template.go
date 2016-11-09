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

func ({{$v}} *{{.Type}}) NewKey(ctx context.Context) *datastore.Key {
    return helper.NewKey(ctx, "{{.Type}}", {{$v}}.{{.IDField}})
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
    noTimestampUpdate bool
}

// Default{{.Type}}Kind is a default value of *{{.Type}}Kind
var Default{{.Type}}Kind = &{{.Type}}Kind{}

const {{.Type}}KindLoggerKey = "ent.{{snakecase .Type}}"

func (k *{{.Type}}Kind) UseDefaultIfNil(b bool) *{{.Type}}Kind {
    k.useDefaultIfNil = b
    return k
}

// Get gets the kind entity from datastore
func (k *{{.Type}}Kind) Get(ctx context.Context, key interface{}) (*datastore.Key, *{{.Type}}, error) {
    keys, ents, err := k.GetMulti(ctx, key)
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

// GetMulti do Get with multiple keys
func (k *{{.Type}}Kind) GetMulti(ctx context.Context, keys ...interface{}) ([]*datastore.Key, []*{{.Type}}, error) {
    var size = len(keys)
    var memKeys []string
    var dsKeys  []*datastore.Key
    var ents []*{{.Type}}
    if size == 0 {
        return nil, nil, nil
    }
    logger := xlog.WithContext(ctx).WithKey({{.Type}}KindLoggerKey)
    dsKeys = make([]*datastore.Key, size, size)
    for i := range keys {
        dsKeys[i] = helper.NewKey(ctx, "{{.Type}}", keys[i])
    }
    ents = make([]*{{.Type}}, size, size)
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
        logger.Debug(func(p *xlog.Printer){
            p.Println("{{.Type}}#GetMulti [Memcache]", )
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
    err := helper.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts)
    if helper.IsDatastoreError(err) {
        // we return nil even some ents hits the cache.
        return nil, nil, err
    }

    if k.useDefaultIfNil {
        for i:=0; i < cacheMissingSize; i++ {
            if cacheMissingEnts[i] == nil {
                cacheMissingEnts[i] = New{{.Type}}()
                cacheMissingEnts[i].{{.IDField}} = dsKeys[i].StringID() // TODO: Support non-string key as ID
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
                logger.Warnf("Failed to create {{.Type}}) caches: %v", err)
            }
        }
    }

    logger.Debug(func(p *xlog.Printer){
        p.Printf(
            "{{.Type}}#GetMulti [Datastore] (UseDefault: %t, NoCache: %t)\n",
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
func (k *{{.Type}}Kind) MustGetMulti(ctx context.Context, keys ...interface{}) []*{{.Type}} {
    _, v, err := k.GetMulti(ctx, keys...)
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
    var dsKeys  []*datastore.Key
    if size == 0 {
        return nil, nil
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

    logger.Debug(func(p *xlog.Printer){
        p.Printf(
            "{{.Type}}#PutMulti [Datastore] (NoCache: %t)\n",
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

// MustPutMulti do Put with multiple keys
func (k *{{.Type}}Kind) MustPutMulti(ctx context.Context, ents []*{{.Type}}) ([]*datastore.Key) {
    keys, err := k.PutMulti(ctx, ents)
    if err != nil {
        panic(err)
    }
    return keys
}

func (k *{{.Type}}Kind) DeleteMulti(ctx context.Context, keys ...interface{}) ([]*datastore.Key, error) {
    var size = len(keys)
    var dsKeys  []*datastore.Key
    if size == 0 {
        return nil, nil
    }
    logger := xlog.WithContext(ctx).WithKey({{.Type}}KindLoggerKey)
    dsKeys = make([]*datastore.Key, size, size)
    for i := range keys {
        dsKeys[i] = helper.NewKey(ctx, "{{.Type}}", keys[i])
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
            memKeys[i] =ent.GetMemcacheKey(dsKeys[i])
        }
        err := memcache.DeleteMulti(ctx, memKeys)
        if memcache.IsMemcacheError(err) {
            logger.Warnf("Failed to invalidate memcache keys: %v", err)
        }
    }

    logger.Debug(func(p *xlog.Printer){
        p.Printf(
            "{{.Type}}#DeleteMulti [Datastore] (NoCache: %t)\n",
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


// {{.Type}}Query helps to build and execute a query
type {{.Type}}Query struct {
    q *helper.Query
}

func New{{.Type}}Query() *{{.Type}}Query {
    return &{{.Type}}Query{
        q: helper.NewQuery("{{.Type}}"),
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
`
