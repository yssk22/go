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

type {{.StructName}}Kind struct {}

func New{{.StructName}}Kind() *{{.StructName}}Kind {
	return {{mkPrivate .StructName}}KindInstance
}

func (d *{{.StructName}}Kind) Get(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, *{{.StructName}}, error) {
    keys, ents, err := d.GetMulti(ctx, []interface{}{key}, options...)
    if err != nil {
        return nil, nil, err
    }
    return keys[0], ents[0], nil
}

func (d *{{.StructName}}Kind) MustGet(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, *{{.StructName}}) {
	k, v, e := d.Get(ctx, key, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{.StructName}}Kind) GetMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, []*{{.StructName}}, error) {
	var err error
	var dsKeys []*datastore.Key
	var ents []*{{.StructName}}
	if dsKeys, err = ds.NormalizeKeys(ctx, "{{.StructName}}", keys); err != nil {
		return nil, nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil, nil
	}
	ents = make([]*{{.StructName}}, size, size)
	if err = ds.GetMulti(ctx, dsKeys, ents, options...); err != nil {
		return nil, nil, err
	}
	return dsKeys, ents, nil
}

func (d *{{.StructName}}Kind) MustGetMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, []*{{.StructName}}) {
	k, v, e := d.GetMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{.StructName}}Kind) Put(ctx context.Context, ent *{{.StructName}}, options ...ds.CRUDOption) (*datastore.Key, error) {
	keys, err := d.PutMulti(ctx, []*{{.StructName}}{ent}, options...)
    if err != nil {
        return nil, err
	}
	return keys[0], nil
}

func (d *{{.StructName}}Kind) MustPut(ctx context.Context, ent *{{.StructName}}, options ...ds.CRUDOption) (*datastore.Key) {
	k, e := d.Put(ctx, ent, options...)
	xerrors.MustNil(e)
	return k
}

func (d *{{.StructName}}Kind) PutMulti(ctx context.Context, ents []*{{.StructName}}, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	dsKeys = make([]*datastore.Key, size, size)
	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
		{{with .TimestampField -}}
		ents[i].{{.}} = xtime.Now()
		{{end -}}
	}
	if dsKeys, err = ds.PutMulti(ctx, dsKeys, ents); err != nil {
		return nil, err
	}
	return dsKeys, nil
}

func (d *{{.StructName}}Kind) MustPutMulti(ctx context.Context, ents []*{{.StructName}}, options ...ds.CRUDOption) ([]*datastore.Key) {
	keys, err := d.PutMulti(ctx, ents, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *{{.StructName}}Kind) Delete(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, error) {
    keys, err := d.DeleteMulti(ctx, []interface{}{key}, options...)
    if err != nil {
        return nil, err
    }
    return keys[0], nil
}

func (d *{{.StructName}}Kind) MustDelete(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key) {
	k, e := d.Delete(ctx, key, options...)
	xerrors.MustNil(e)
	return k
}

func (d *{{.StructName}}Kind) DeleteMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	var err error
	var dsKeys []*datastore.Key
	if dsKeys, err = ds.NormalizeKeys(ctx, "{{.StructName}}", keys); err != nil {
		return nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil
	}
	if err = ds.DeleteMulti(ctx, dsKeys); err != nil  {
		return nil, xerrors.Wrap(err, "datastore error")
	}
	return dsKeys, nil
}

func (d *{{.StructName}}Kind) MustDeleteMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key) {
	k, e := d.DeleteMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k
}

func (d *{{.StructName}}Kind) DeleteMatched(ctx context.Context, q *{{.StructName}}Query, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	keys, err := q.query.KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	_, err = d.DeleteMulti(ctx, keys, options...)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (d *{{.StructName}}Kind) MustDeleteMatched(ctx context.Context, q *{{.StructName}}Query, options ...ds.CRUDOption) ([]*datastore.Key) {
	keys, err := d.DeleteMatched(ctx, q, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *{{.StructName}}Kind) Replace(ctx context.Context, ent *{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.CRUDOption) (*datastore.Key, *{{.StructName}}, error) {
    keys, ents, err := d.ReplaceMulti(ctx, []*{{.StructName}}{ent}, replacer, options...)
    if err != nil {
        return nil, ents[0], err
	}
    return keys[0], ents[0], err
}

func (d *{{.StructName}}Kind) MustReplace(ctx context.Context, ent *{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.CRUDOption) (*datastore.Key, *{{.StructName}}) {
	k, v, e := d.Replace(ctx, ent, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{.StructName}}Kind) ReplaceMulti(ctx context.Context, ents []*{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.CRUDOption) ([]*datastore.Key, []*{{.StructName}}, error) {
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

func (d *{{.StructName}}Kind) MustReplaceMulti(ctx context.Context, ents []*{{.StructName}}, replacer {{.StructName}}Replacer, options ...ds.CRUDOption) ([]*datastore.Key, []*{{.StructName}}) {
	k, v, e := d.ReplaceMulti(ctx, ents, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

type {{.StructName}}Query struct {
	query *ds.Query
	viaKeys bool
}

func New{{.StructName}}Query() *{{.StructName}}Query {
	return &{{.StructName}}Query{
		query: ds.NewQuery("{{.KindName}}"),
		viaKeys: false,
	}
}

{{queryFuncs .}}

func (d *{{.StructName}}Query) Start(s string) *{{.StructName}}Query {
	d.query = d.query.Start(s)
	return d
}

func (d *{{.StructName}}Query) End(s string) *{{.StructName}}Query {
	d.query = d.query.End(s)
	return d
}

func (d *{{.StructName}}Query) Limit(n int) *{{.StructName}}Query {
	d.query = d.query.Limit(n)
	return d
}

func (d *{{.StructName}}Query) ViaKeys() *{{.StructName}}Query {
	d.viaKeys = true
	return d
}

func (d *{{.StructName}}Query) GetAll(ctx context.Context) ([]*datastore.Key, []{{.StructName}}, error) {
	if d.viaKeys {
		keys, err := d.query.KeysOnly().GetAll(ctx, nil)
		if err != nil {
			return nil, nil, err
		}
		_, ents, err := {{snakecase .StructName}}KindInstance.GetMulti(ctx, keys)
		if err != nil {
			return nil, nil, err
		}
		list := make([]{{.StructName}}, len(ents))
		for i, e := range ents {
			list[i] = *e
		}
		return keys, list, nil
	}
	var ent []{{.StructName}}
	keys, err := d.query.GetAll(ctx, &ent)
	if err != nil {
		return nil, nil, err
	}
	return keys, ent, nil
}

func (d *{{.StructName}}Query) MustGetAll(ctx context.Context) ([]*datastore.Key, []{{.StructName}}) {
	keys, ents, err := d.GetAll(ctx)
	xerrors.MustNil(err)
	return keys, ents
}

func (d *{{.StructName}}Query) Count(ctx context.Context) (int, error) {
	return d.query.Count(ctx)
}

func (d *{{.StructName}}Query) MustCount(ctx context.Context) (int) {
	c, err := d.query.Count(ctx)
	xerrors.MustNil(err)
	return c
}

func (d *{{.StructName}}Query) Run(ctx context.Context) (*{{.StructName}}Iterator, error) {
	iter, err := d.query.Run(ctx)
	if err != nil {
		return nil, err
	}
	return &{{.StructName}}Iterator{
		ctx: ctx,
		iter: iter,
		viaKeys: d.viaKeys,
	}, err
}

func (d *{{.StructName}}Query) MustRun(ctx context.Context) (*{{.StructName}}Iterator) {
	iter, err := d.Run(ctx)
	xerrors.MustNil(err)
	return iter
}

func (d *{{.StructName}}Query) RunAll(ctx context.Context) ([]datastore.Key, []{{.StructName}}, string, error) {
	iter, err := d.Run(ctx)
	if err != nil {
		return nil, nil, "", err
	}
	var keys []datastore.Key
	var ents []{{.StructName}}
	for {
		key, ent, err := iter.Next()
		if err != nil{
			return nil, nil, "", err
		}
		if ent == nil {
			cursor, err := iter.iter.Cursor()
			if err != nil {
				return nil, nil, "", err
			}
			return keys, ents, cursor.String(), nil
		}
		keys = append(keys, *key)
		ents = append(ents, *ent)
	}
}

func (d *{{.StructName}}Query) MustRunAll(ctx context.Context) ([]datastore.Key, []{{.StructName}}, string) {
	keys, ents, next, err := d.RunAll(ctx)
	xerrors.MustNil(err)
	return keys, ents, next
}

type {{.StructName}}Iterator struct {
	ctx context.Context
	iter *datastore.Iterator
	viaKeys bool
}

func (iter *{{.StructName}}Iterator) Cursor() (datastore.Cursor, error) {
	return iter.iter.Cursor()
}

func (iter *{{.StructName}}Iterator) MustCursor() (datastore.Cursor) {
	c, err := iter.iter.Cursor()
	xerrors.MustNil(err)
	return c
}

func (iter *{{.StructName}}Iterator) Next() (*datastore.Key, *{{.StructName}}, error) {
	if iter.viaKeys {
		key, err := iter.iter.Next(nil)
		if err != nil {
			if err == datastore.Done {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		_, ent, err := {{snakecase .StructName}}KindInstance.Get(iter.ctx, key)
		if err != nil {
			return nil, nil, err
		}
		return key, ent, nil

	}
	var ent {{.StructName}}
	key, err := iter.iter.Next(&ent)
	if err != nil {
		if err == datastore.Done {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return key, &ent, nil
}

func (iter *{{.StructName}}Iterator) MustNext() (*datastore.Key, *{{.StructName}}) {
	key, ent, err := iter.Next()
	xerrors.MustNil(err)
	return key, ent
}

var {{mkPrivate .StructName}}KindInstance = &{{.StructName}}Kind{}

{{end -}}
`
