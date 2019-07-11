package typed

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

type {{.StructName}}KindClient struct {
	client *ds.Client
}

func New{{.StructName}}KindClient(client *ds.Client) *{{.StructName}}KindClient {
	return &{{.StructName}}KindClient{
		client: client,
	}
}

func (d *{{.StructName}}KindClient) Get(ctx context.Context, key interface{}) (*datastore.Key, *{{.StructName}}, error) {
    keys, ents, err := d.GetMulti(ctx, []interface{}{key})
    if err != nil {
        return nil, nil, err
    }
    return keys[0], ents[0], nil
}

func (d *{{.StructName}}KindClient) MustGet(ctx context.Context, key interface{}) (*datastore.Key, *{{.StructName}}) {
	k, v, e := d.Get(ctx, key)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{.StructName}}KindClient) GetMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, []*{{.StructName}}, error) {
	var err error
	var dsKeys []*datastore.Key
	var ents []*{{.StructName}}
	if dsKeys, err = ds.NormalizeKeys(ctx, "{{.KindName}}", keys); err != nil {
		return nil, nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil, nil
	}
	ents = make([]*{{.StructName}}, size, size)
	if err = d.client.GetMulti(ctx, dsKeys, ents); err != nil {
		return nil, nil, err
	}
	return dsKeys, ents, nil
}

func (d *{{.StructName}}KindClient) MustGetMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, []*{{.StructName}}) {
	k, v, e := d.GetMulti(ctx, keys)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{.StructName}}KindClient) Put(ctx context.Context, ent *{{.StructName}}) (*datastore.Key, error) {
	keys, err := d.PutMulti(ctx, []*{{.StructName}}{ent})
    if err != nil {
        return nil, err
	}
	return keys[0], nil
}

func (d *{{.StructName}}KindClient) MustPut(ctx context.Context, ent *{{.StructName}}) (*datastore.Key) {
	k, e := d.Put(ctx, ent)
	xerrors.MustNil(e)
	return k
}

func (d *{{.StructName}}KindClient) PutMulti(ctx context.Context, ents []*{{.StructName}}) ([]*datastore.Key, error) {
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	dsKeys = make([]*datastore.Key, size, size)
	if size == 0 {
		return nil, nil
	}
	_, hasBeforeSave := interface{}(ents[0]).(ds.BeforeSave)
	_, hasAfterSave := interface{}(ents[0]).(ds.AfterSave)

	if hasBeforeSave {
		for i := range ents {
			if err := interface{}(ents[i]).(ds.BeforeSave).BeforeSave(ctx); err != nil {
				return nil, err
			}
		}
	}

	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
		{{with .TimestampField -}}
		ents[i].{{.}} = xtime.Now()
		{{end -}}
	}
	if dsKeys, err = d.client.PutMulti(ctx, dsKeys, ents); err != nil {
		return nil, err
	}

	if hasAfterSave {
		for i := range ents {
			if err := interface{}(ents[i]).(ds.AfterSave).AfterSave(ctx); err != nil {
				return nil, err
			}
		}
	}
	return dsKeys, nil
}

func (d *{{.StructName}}KindClient) MustPutMulti(ctx context.Context, ents []*{{.StructName}}) ([]*datastore.Key) {
	keys, err := d.PutMulti(ctx, ents)
	xerrors.MustNil(err)
	return keys
}

func (d *{{.StructName}}KindClient) Delete(ctx context.Context, key interface{}) (*datastore.Key, error) {
    keys, err := d.DeleteMulti(ctx, []interface{}{key})
    if err != nil {
        return nil, err
    }
    return keys[0], nil
}

func (d *{{.StructName}}KindClient) MustDelete(ctx context.Context, key interface{}) (*datastore.Key) {
	k, e := d.Delete(ctx, key)
	xerrors.MustNil(e)
	return k
}

func (d *{{.StructName}}KindClient) DeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key, error) {
	var err error
	var dsKeys []*datastore.Key
	if dsKeys, err = ds.NormalizeKeys(ctx, "{{.KindName}}", keys); err != nil {
		return nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil
	}
	if err = d.client.DeleteMulti(ctx, dsKeys); err != nil  {
		return nil, xerrors.Wrap(err, "datastore error")
	}
	return dsKeys, nil
}

func (d *{{.StructName}}KindClient) MustDeleteMulti(ctx context.Context, keys interface{}) ([]*datastore.Key) {
	k, e := d.DeleteMulti(ctx, keys)
	xerrors.MustNil(e)
	return k
}

func (d *{{.StructName}}KindClient) DeleteMatched(ctx context.Context, q *{{.StructName}}Query) ([]*datastore.Key, error) {
	keys, err := d.client.GetAll(ctx, q.query.KeysOnly(), nil)
	if err != nil {
		return nil, err
	}
	_, err = d.DeleteMulti(ctx, keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (d *{{.StructName}}KindClient) MustDeleteMatched(ctx context.Context, q *{{.StructName}}Query) ([]*datastore.Key) {
	keys, err := d.DeleteMatched(ctx, q)
	xerrors.MustNil(err)
	return keys
}

func (d *{{.StructName}}KindClient) Replace(ctx context.Context, ent *{{.StructName}}, replacer {{.StructName}}Replacer) (*datastore.Key, *{{.StructName}}, error) {
    keys, ents, err := d.ReplaceMulti(ctx, []*{{.StructName}}{ent}, replacer)
    if err != nil {
        return nil, ents[0], err
	}
    return keys[0], ents[0], err
}

func (d *{{.StructName}}KindClient) MustReplace(ctx context.Context, ent *{{.StructName}}, replacer {{.StructName}}Replacer) (*datastore.Key, *{{.StructName}}) {
	k, v, e := d.Replace(ctx, ent, replacer)
	xerrors.MustNil(e)
	return k, v
}

func (d *{{.StructName}}KindClient) ReplaceMulti(ctx context.Context, ents []*{{.StructName}}, replacer {{.StructName}}Replacer) ([]*datastore.Key, []*{{.StructName}}, error) {
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

func (d *{{.StructName}}KindClient) MustReplaceMulti(ctx context.Context, ents []*{{.StructName}}, replacer {{.StructName}}Replacer) ([]*datastore.Key, []*{{.StructName}}) {
	k, v, e := d.ReplaceMulti(ctx, ents, replacer)
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

func (q *{{.StructName}}Query) Start(s string) *{{.StructName}}Query {
	q.query = q.query.Start(s)
	return q
}

func (q *{{.StructName}}Query) End(s string) *{{.StructName}}Query {
	q.query = q.query.End(s)
	return q
}

func (q *{{.StructName}}Query) Limit(n int) *{{.StructName}}Query {
	q.query = q.query.Limit(n)
	return q
}

func (q *{{.StructName}}Query) ViaKeys() *{{.StructName}}Query {
	q.viaKeys = true
	return q
}

func (d *{{.StructName}}KindClient) GetAll(ctx context.Context, q *{{.StructName}}Query) ([]*datastore.Key, []{{.StructName}}, error) {
	if q.viaKeys {
		keys, err := d.client.GetAll(ctx, q.query.KeysOnly(), nil)
		if err != nil {
			return nil, nil, err
		}
		ents := make([]*{{.StructName}}, len(keys))
		err = d.client.GetMulti(ctx, keys, ents)
		if err != nil {
			return nil, nil, err
		}
		result := make([]{{.StructName}}, 0)
		for _, e := range ents {
			if e != nil {
				result = append(result, *e)
			}
		}
		return keys, result, nil
	} else {
		var ent []{{.StructName}}
		keys, err := d.client.GetAll(ctx, q.query, &ent)
		if err != nil {
			return nil, nil, err
		}
		return keys, ent, nil
	}
}

func (d *{{.StructName}}KindClient) GetOne(ctx context.Context, q *{{.StructName}}Query) (*datastore.Key, *{{.StructName}}, error) {
	keys, ents, err := d.GetAll(ctx, q.Limit(1))
	if err != nil {
		return nil, nil, err
	}
	if len(keys) == 0 {
		return nil, nil, nil
	}
	return keys[0], &(ents[0]), nil
}

func (d *{{.StructName}}KindClient) MustGetAll(ctx context.Context, q *{{.StructName}}Query) ([]*datastore.Key, []{{.StructName}}) {
	keys, ents, err := d.GetAll(ctx, q)
	xerrors.MustNil(err)
	return keys, ents
}

func (d *{{.StructName}}KindClient) Count(ctx context.Context, q *{{.StructName}}Query) (int, error) {
	return d.client.Count(ctx, q.query)
}

func (d *{{.StructName}}KindClient) MustCount(ctx context.Context,  q *{{.StructName}}Query) (int) {
	c, err := d.Count(ctx, q)
	xerrors.MustNil(err)
	return c
}

func (d *{{.StructName}}KindClient) Run(ctx context.Context, q *{{.StructName}}Query) (*{{.StructName}}Iterator, error) {
	iter, err := d.client.Run(ctx, q.query)
	if err != nil {
		return nil, err
	}
	return &{{.StructName}}Iterator{
		ctx: ctx,
		iter: iter,
		viaKeys: q.viaKeys,
		client: d,
	}, err
}

func (d *{{.StructName}}KindClient) MustRun(ctx context.Context, q *{{.StructName}}Query) (*{{.StructName}}Iterator) {
	iter, err := d.Run(ctx, q)
	xerrors.MustNil(err)
	return iter
}

func (d *{{.StructName}}KindClient) RunAll(ctx context.Context, q *{{.StructName}}Query) ([]datastore.Key, []{{.StructName}}, string, error) {
	iter, err := d.Run(ctx, q)
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

func (d *{{.StructName}}KindClient) MustRunAll(ctx context.Context, q *{{.StructName}}Query) ([]datastore.Key, []{{.StructName}}, string) {
	keys, ents, next, err := d.RunAll(ctx, q)
	xerrors.MustNil(err)
	return keys, ents, next
}

type {{.StructName}}Iterator struct {
	ctx context.Context
	iter *datastore.Iterator
	viaKeys bool
	client *{{.StructName}}KindClient
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
			if err == iterator.Done {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		_, ent, err := iter.client.Get(iter.ctx, key)
		if err != nil {
			return nil, nil, err
		}
		return key, ent, nil
	}
	var ent {{.StructName}}
	key, err := iter.iter.Next(&ent)
	if err != nil {
		if err == iterator.Done {
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

{{end -}}
`
