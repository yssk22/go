package datastore

import "github.com/yssk22/go/generator"

type bindings struct {
	Package    string
	Dependency *generator.Dependency
	Specs      []*Spec
}

const templateFile = `
package {{.Package}}

{{.Dependency.GenImport}}

{{range .Specs -}}

func (s *{{.StructName}}) NewKey(ctx context.Context) *datastore.Key {
	return ds.NewKey(ctx, "{{.KindName}}", s.{{.KeyField}})
}

type {{.StructName}}Replacer interface {
	Replace(*{{.StructName}}, *{{.StructName}}) *{{.StructName}}
}

type {{.StructName}}ReplacerFunc func(*{{.StructName}}, *{{.StructName}}) *{{.StructName}}

func (f {{.StructName}}ReplacerFunc) Replace(old *{{.StructName}}, new *{{.StructName}}) *{{.StructName}} {
	return f(old, new)
}

type {{.StructName}}Datastore interface {
	Get(context.Context, interface{}, ...ds.Option) (*datastore.Key, *{{.StructName}}, error)
	MustGet(context.Context, interface{}, ...ds.Option) (*datastore.Key, *{{.StructName}})
	GetMulti(context.Context, interface{}, ...ds.Option) ([]*datastore.Key, []*{{.StructName}}, error)
	MustGetMulti(context.Context, interface{}, ...ds.Option) ([]*datastore.Key, []*{{.StructName}})

	Put(context.Context, *{{.StructName}}, ...ds.Option) (*datastore.Key, error)
	MustPut(context.Context, *{{.StructName}}, ...ds.Option) (*datastore.Key)
	PutMulti(context.Context, []*{{.StructName}}, ...ds.Option) ([]*datastore.Key, error)
	MustPutMulti(context.Context, []*{{.StructName}}, ...ds.Option) ([]*datastore.Key)

	Delete(context.Context, interface{}, ...ds.Option) (*datastore.Key, error)
	MustDelete(context.Context, interface{}, ...ds.Option) (*datastore.Key)
	DeleteMulti(context.Context, interface{}, ...ds.Option) ([]*datastore.Key, error)
	MustDeleteMulti(context.Context, interface{}, ...ds.Option) ([]*datastore.Key)

	Replace(context.Context, *{{.StructName}}, {{.StructName}}Replacer, ...ds.Option) (*datastore.Key, *{{.StructName}}, error)
	MustReplace(context.Context, *{{.StructName}}, {{.StructName}}Replacer, ...ds.Option) (*datastore.Key, *{{.StructName}})
	ReplaceMulti(context.Context, []*{{.StructName}}, {{.StructName}}Replacer, ...ds.Option) ([]*datastore.Key, []*{{.StructName}}, error)
	MustReplaceMulti(context.Context, []*{{.StructName}}, {{.StructName}}Replacer, ...ds.Option) ([]*datastore.Key, []*{{.StructName}})
}

func Get{{.StructName}}Datastore() {{.StructName}}Datastore {
	return {{snakecase .StructName}}Instance
}

type {{snakecase .StructName}}Datastore struct {}

func (d *{{snakecase .StructName}}Datastore) Get(ctx context.Context, key interface{}, options ...ds.Option) (*datastore.Key, *{{.StructName}}, error) {
    keys, ents, err := d.GetMulti(ctx, []interface{}{key}, options...)
    if err != nil {
        return nil, nil, err
    }
    return keys[0], ents[0], nil
}

func (d *{{snakecase .StructName}}Datastore) MustGet(ctx context.Context, key interface{}, options ...ds.Option) (*datastore.Key, *{{.StructName}}) {
	k, v, e := d.Get(ctx, key, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{snakecase .StructName}}Datastore) GetMulti(ctx context.Context, keys interface{}, options ...ds.Option) ([]*datastore.Key, []*{{.StructName}}, error) {
	var err error
	var dsKeys []*datastore.Key
	var memKeys []string
	var ents []*{{.StructName}}

	opts := ds.NewCRUDOption(options...)
	if opts.Namespace != nil {
		ctx, err = appengine.Namespace(ctx, *(opts.Namespace))
		if err != nil {
			return nil, nil, xerrors.Wrap(err, "cannot enforce namespace")
		}
	}

	if dsKeys, err = ds.NormalizeKeys(ctx, "{{.StructName}}", keys); err != nil {
		return nil, nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}

	size := len(dsKeys)
	if size == 0 {
		return nil, nil, nil
	}

	ents = make([]*{{.StructName}}, size, size)
	// fetch from cache
	if !opts.NoCache {
		memKeys = make([]string, size, size)
		for i := range dsKeys {
			memKeys[i] = ds.GetMemcacheKey(dsKeys[i])
		}
		err = memcache.GetMulti(ctx, memKeys, ents)
		if err == nil { 
			// Hit caches on all keys!!
			return dsKeys, ents, nil
		}
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
	cacheMissingEnts := make([]*{{.StructName}}, cacheMissingSize, cacheMissingSize)	
	if err = ds.GetMulti(ctx, cacheMissingKeys, cacheMissingEnts); ds.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, nil, xerrors.Wrap(err, "datastore error")
	}
	for i := range cacheMissingKeys {
		entIdx := key2Idx[cacheMissingKeys[i]]
		ents[entIdx] = cacheMissingEnts[i]
	}

	// udpate cache
	if !opts.NoCache {
		cacheEnts := make([]*{{.StructName}}, 0)
		cacheKeys := make([]string, 0)
		for i := range ents {
			if ents[i] != nil {
				cacheEnts = append(cacheEnts, ents[i])
				cacheKeys = append(cacheKeys, memKeys[i])
			}
		}
		if len(cacheEnts) > 0 {
			if err := memcache.SetMulti(ctx, cacheKeys, cacheEnts); err != nil {
			}
		}
	}

	return dsKeys, ents, nil
}

func (d *{{snakecase .StructName}}Datastore) MustGetMulti(ctx context.Context, keys interface{}, options ...ds.Option) ([]*datastore.Key, []*{{.StructName}}) {
	k, v, e := d.GetMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{snakecase .StructName}}Datastore) Put(ctx context.Context, ent *{{.StructName}}, options ...ds.Option) (*datastore.Key, error) {
	keys, err := d.PutMulti(ctx, []*{{.StructName}}{ent}, options...)
    if err != nil {
        return nil, err
	}
	return keys[0], nil
}

func (d *{{snakecase .StructName}}Datastore) MustPut(ctx context.Context, ent *{{.StructName}}, options ...ds.Option) (*datastore.Key) {
	k, e := d.Put(ctx, ent, options...)
	xerrors.MustNil(e)
	return k
}

func (d *{{snakecase .StructName}}Datastore) PutMulti(ctx context.Context, ents []*{{.StructName}}, options ...ds.Option) ([]*datastore.Key, error) {
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	if size == 0 {
		return nil, nil
	}
	if size >= ds.MaxEntitiesPerUpdate {
		return nil, ds.ErrTooManyEnts
	}
	opts := ds.NewCRUDOption(options...)
	if opts.Namespace != nil {
		ctx, err = appengine.Namespace(ctx, *(opts.Namespace))
		if err != nil {
			return nil, xerrors.Wrap(err, "cannot enforce namespace")
		}
	}

	dsKeys = make([]*datastore.Key, size, size)
	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
	}

	if !opts.NoTimestampUpdate {
		for i := range ents {
			ents[i].{{.TimestampField}} = xtime.Now()
		}
	}

	if _, err = ds.PutMulti(ctx, dsKeys, ents); ds.IsDatastoreError(err) {
		return nil, xerrors.Wrap(err, "datastore error")
	}

	if !opts.NoCache {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = ds.GetMemcacheKey(dsKeys[i])
		}
		if err := memcache.DeleteMulti(ctx, memKeys); memcache.IsMemcacheError(err) {
		}
	}

	return dsKeys, nil
}

func (d *{{snakecase .StructName}}Datastore) MustPutMulti(ctx context.Context, ents []*{{.StructName}}, options ...ds.Option) ([]*datastore.Key) {
	keys, err := d.PutMulti(ctx, ents, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *{{snakecase .StructName}}Datastore) Delete(ctx context.Context, key interface{}, options ...ds.Option) (*datastore.Key, error) {
    keys, err := d.DeleteMulti(ctx, []interface{}{key}, options...)
    if err != nil {
        return nil, err
    }
    return keys[0], nil
}

func (d *{{snakecase .StructName}}Datastore) MustDelete(ctx context.Context, key interface{}, options ...ds.Option) (*datastore.Key) {
	k, e := d.Delete(ctx, key, options...)
	xerrors.MustNil(e)
	return k
}

func (d *{{snakecase .StructName}}Datastore) DeleteMulti(ctx context.Context, keys interface{}, options ...ds.Option) ([]*datastore.Key, error) {
	var err error
	var dsKeys []*datastore.Key

	opts := ds.NewCRUDOption(options...)
	if opts.Namespace != nil {
		ctx, err = appengine.Namespace(ctx, *(opts.Namespace))
		if err != nil {
			return nil, xerrors.Wrap(err, "cannot enforce namespace")
		}
	}

	if dsKeys, err = ds.NormalizeKeys(ctx, "{{.StructName}}", keys); err != nil {
		return nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}

	size := len(dsKeys)
	if size == 0 {
		return nil, nil
	}
	if size >= ds.MaxEntitiesPerUpdate {
		return nil, ds.ErrTooManyEnts
	}

	if err = ds.DeleteMulti(ctx, dsKeys); ds.IsDatastoreError(err) {
		// we return nil even some ents hits the cache.
		return nil, xerrors.Wrap(err, "datastore error")
	}

	// invalidate cache
	if !opts.NoCache {
		memKeys := make([]string, size, size)
		for i := range memKeys {
			memKeys[i] = ds.GetMemcacheKey(dsKeys[i])
		}
		if err = memcache.DeleteMulti(ctx, memKeys); memcache.IsMemcacheError(err) {
		}
	}

	return dsKeys, nil
}

func (d *{{snakecase .StructName}}Datastore) MustDeleteMulti(ctx context.Context, keys interface{}, options ...ds.Option) ([]*datastore.Key) {
	k, e := d.DeleteMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k
}

func (d *{{snakecase .StructName}}Datastore) Replace(ctx context.Context, ent *{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.Option) (*datastore.Key, *{{.StructName}}, error) {
    keys, ents, err := d.ReplaceMulti(ctx, []*{{.StructName}}{ent}, replacer, options...)
    if err != nil {
        return nil, ents[0], err
	}
    return keys[0], ents[0], err
}

func (d *{{snakecase .StructName}}Datastore) MustReplace(ctx context.Context, ent *{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.Option) (*datastore.Key, *{{.StructName}}) {
	k, v, e := d.Replace(ctx, ent, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{snakecase .StructName}}Datastore) ReplaceMulti(ctx context.Context, ents []*{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.Option) ([]*datastore.Key, []*{{.StructName}}, error) {
	var size = len(ents)
	var dsKeys = make([]*datastore.Key, size, size)
	if size == 0 {
		return dsKeys, ents, nil
	}
	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
	}
	_, existing, err := d.GetMulti(ctx, dsKeys)
	if err != nil {
		return nil, ents, err
	}
	for i, exist := range existing {
		if exist != nil {
			ents[i] = replacer.Replace(exist, ents[i])
		}
	}
	dsKeys, err = d.PutMulti(ctx, ents)
	return dsKeys, ents, err
}

func (d *{{snakecase .StructName}}Datastore) MustReplaceMulti(ctx context.Context, ents []*{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.Option) ([]*datastore.Key, []*{{.StructName}}) {
	k, v, e := d.ReplaceMulti(ctx, ents, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

var {{snakecase .StructName}}Instance = &{{snakecase .StructName}}Datastore{}

{{end -}}
`
