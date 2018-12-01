// Code generated by github.com/yssk22/go/generator DO NOT EDIT.
//
package config

import (
	"context"
	ds "github.com/yssk22/go/gae/datastore"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xtime"
	"google.golang.org/appengine/datastore"
)

func (s *ServiceConfig) NewKey(ctx context.Context) *datastore.Key {
	return ds.NewKey(ctx, "ServiceConfig", s.Key)
}

type ServiceConfigReplacer interface {
	Replace(*ServiceConfig, *ServiceConfig) *ServiceConfig
}

type ServiceConfigReplacerFunc func(*ServiceConfig, *ServiceConfig) *ServiceConfig

func (f ServiceConfigReplacerFunc) Replace(old *ServiceConfig, new *ServiceConfig) *ServiceConfig {
	return f(old, new)
}

type ServiceConfigKind struct{}

func NewServiceConfigKind() *ServiceConfigKind {
	return serviceConfigKindInstance
}

func (d *ServiceConfigKind) Get(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, *ServiceConfig, error) {
	keys, ents, err := d.GetMulti(ctx, []interface{}{key}, options...)
	if err != nil {
		return nil, nil, err
	}
	return keys[0], ents[0], nil
}

func (d *ServiceConfigKind) MustGet(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, *ServiceConfig) {
	k, v, e := d.Get(ctx, key, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *ServiceConfigKind) GetMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, []*ServiceConfig, error) {
	var err error
	var dsKeys []*datastore.Key
	var ents []*ServiceConfig
	if dsKeys, err = ds.NormalizeKeys(ctx, "ServiceConfig", keys); err != nil {
		return nil, nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil, nil
	}
	ents = make([]*ServiceConfig, size, size)
	if err = ds.GetMulti(ctx, dsKeys, ents, options...); err != nil {
		return nil, nil, err
	}
	return dsKeys, ents, nil
}

func (d *ServiceConfigKind) MustGetMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, []*ServiceConfig) {
	k, v, e := d.GetMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *ServiceConfigKind) Put(ctx context.Context, ent *ServiceConfig, options ...ds.CRUDOption) (*datastore.Key, error) {
	keys, err := d.PutMulti(ctx, []*ServiceConfig{ent}, options...)
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

func (d *ServiceConfigKind) MustPut(ctx context.Context, ent *ServiceConfig, options ...ds.CRUDOption) *datastore.Key {
	k, e := d.Put(ctx, ent, options...)
	xerrors.MustNil(e)
	return k
}

func (d *ServiceConfigKind) PutMulti(ctx context.Context, ents []*ServiceConfig, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	dsKeys = make([]*datastore.Key, size, size)
	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
		ents[i].UpdatedAt = xtime.Now()
	}
	if dsKeys, err = ds.PutMulti(ctx, dsKeys, ents); err != nil {
		return nil, err
	}
	return dsKeys, nil
}

func (d *ServiceConfigKind) MustPutMulti(ctx context.Context, ents []*ServiceConfig, options ...ds.CRUDOption) []*datastore.Key {
	keys, err := d.PutMulti(ctx, ents, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *ServiceConfigKind) Delete(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, error) {
	keys, err := d.DeleteMulti(ctx, []interface{}{key}, options...)
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

func (d *ServiceConfigKind) MustDelete(ctx context.Context, key interface{}, options ...ds.CRUDOption) *datastore.Key {
	k, e := d.Delete(ctx, key, options...)
	xerrors.MustNil(e)
	return k
}

func (d *ServiceConfigKind) DeleteMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	var err error
	var dsKeys []*datastore.Key
	if dsKeys, err = ds.NormalizeKeys(ctx, "ServiceConfig", keys); err != nil {
		return nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil
	}
	if err = ds.DeleteMulti(ctx, dsKeys); err != nil {
		return nil, xerrors.Wrap(err, "datastore error")
	}
	return dsKeys, nil
}

func (d *ServiceConfigKind) MustDeleteMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) []*datastore.Key {
	k, e := d.DeleteMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k
}

func (d *ServiceConfigKind) DeleteMatched(ctx context.Context, q *ServiceConfigQuery, options ...ds.CRUDOption) ([]*datastore.Key, error) {
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

func (d *ServiceConfigKind) MustDeleteMatched(ctx context.Context, q *ServiceConfigQuery, options ...ds.CRUDOption) []*datastore.Key {
	keys, err := d.DeleteMatched(ctx, q, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *ServiceConfigKind) Replace(ctx context.Context, ent *ServiceConfig, replacer ServiceConfigReplacer, options ...ds.CRUDOption) (*datastore.Key, *ServiceConfig, error) {
	keys, ents, err := d.ReplaceMulti(ctx, []*ServiceConfig{ent}, replacer, options...)
	if err != nil {
		return nil, ents[0], err
	}
	return keys[0], ents[0], err
}

func (d *ServiceConfigKind) MustReplace(ctx context.Context, ent *ServiceConfig, replacer ServiceConfigReplacer, options ...ds.CRUDOption) (*datastore.Key, *ServiceConfig) {
	k, v, e := d.Replace(ctx, ent, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *ServiceConfigKind) ReplaceMulti(ctx context.Context, ents []*ServiceConfig, replacer ServiceConfigReplacer, options ...ds.CRUDOption) ([]*datastore.Key, []*ServiceConfig, error) {
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

func (d *ServiceConfigKind) MustReplaceMulti(ctx context.Context, ents []*ServiceConfig, replacer ServiceConfigReplacer, options ...ds.CRUDOption) ([]*datastore.Key, []*ServiceConfig) {
	k, v, e := d.ReplaceMulti(ctx, ents, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

type ServiceConfigQuery struct {
	query   *ds.Query
	viaKeys bool
}

func NewServiceConfigQuery() *ServiceConfigQuery {
	return &ServiceConfigQuery{
		query:   ds.NewQuery("ServiceConfig"),
		viaKeys: false,
	}
}

func (d *ServiceConfigQuery) EqKey(v string) *ServiceConfigQuery {
	d.query = d.query.Eq("Key", v)
	return d
}

func (d *ServiceConfigQuery) LtKey(v string) *ServiceConfigQuery {
	d.query = d.query.Lt("Key", v)
	return d
}

func (d *ServiceConfigQuery) LeKey(v string) *ServiceConfigQuery {
	d.query = d.query.Le("Key", v)
	return d
}

func (d *ServiceConfigQuery) GtKey(v string) *ServiceConfigQuery {
	d.query = d.query.Gt("Key", v)
	return d
}

func (d *ServiceConfigQuery) GeKey(v string) *ServiceConfigQuery {
	d.query = d.query.Ge("Key", v)
	return d
}

func (d *ServiceConfigQuery) NeKey(v string) *ServiceConfigQuery {
	d.query = d.query.Ne("Key", v)
	return d
}

func (d *ServiceConfigQuery) AscKey() *ServiceConfigQuery {
	d.query = d.query.Asc("Key")
	return d
}

func (d *ServiceConfigQuery) DescKey() *ServiceConfigQuery {
	d.query = d.query.Desc("Key")
	return d
}

func (d *ServiceConfigQuery) Start(s string) *ServiceConfigQuery {
	d.query = d.query.Start(s)
	return d
}

func (d *ServiceConfigQuery) End(s string) *ServiceConfigQuery {
	d.query = d.query.End(s)
	return d
}

func (d *ServiceConfigQuery) Limit(n int) *ServiceConfigQuery {
	d.query = d.query.Limit(n)
	return d
}

func (d *ServiceConfigQuery) ViaKeys() *ServiceConfigQuery {
	d.viaKeys = true
	return d
}

func (d *ServiceConfigQuery) GetAll(ctx context.Context) ([]*datastore.Key, []ServiceConfig, error) {
	if d.viaKeys {
		keys, err := d.query.KeysOnly().GetAll(ctx, nil)
		if err != nil {
			return nil, nil, err
		}
		_, ents, err := serviceConfigKindInstance.GetMulti(ctx, keys)
		if err != nil {
			return nil, nil, err
		}
		list := make([]ServiceConfig, len(ents))
		for i, e := range ents {
			list[i] = *e
		}
		return keys, list, nil
	}
	var ent []ServiceConfig
	keys, err := d.query.GetAll(ctx, &ent)
	if err != nil {
		return nil, nil, err
	}
	return keys, ent, nil
}

func (d *ServiceConfigQuery) MustGetAll(ctx context.Context) ([]*datastore.Key, []ServiceConfig) {
	keys, ents, err := d.GetAll(ctx)
	xerrors.MustNil(err)
	return keys, ents
}

func (d *ServiceConfigQuery) Count(ctx context.Context) (int, error) {
	return d.query.Count(ctx)
}

func (d *ServiceConfigQuery) MustCount(ctx context.Context) int {
	c, err := d.query.Count(ctx)
	xerrors.MustNil(err)
	return c
}

func (d *ServiceConfigQuery) Run(ctx context.Context) (*ServiceConfigIterator, error) {
	iter, err := d.query.Run(ctx)
	if err != nil {
		return nil, err
	}
	return &ServiceConfigIterator{
		ctx:     ctx,
		iter:    iter,
		viaKeys: d.viaKeys,
	}, err
}

func (d *ServiceConfigQuery) MustRun(ctx context.Context) *ServiceConfigIterator {
	iter, err := d.Run(ctx)
	xerrors.MustNil(err)
	return iter
}

func (d *ServiceConfigQuery) RunAll(ctx context.Context) ([]datastore.Key, []ServiceConfig, string, error) {
	iter, err := d.Run(ctx)
	if err != nil {
		return nil, nil, "", err
	}
	var keys []datastore.Key
	var ents []ServiceConfig
	for {
		key, ent, err := iter.Next()
		if err != nil {
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

func (d *ServiceConfigQuery) MustRunAll(ctx context.Context) ([]datastore.Key, []ServiceConfig, string) {
	keys, ents, next, err := d.RunAll(ctx)
	xerrors.MustNil(err)
	return keys, ents, next
}

type ServiceConfigIterator struct {
	ctx     context.Context
	iter    *datastore.Iterator
	viaKeys bool
}

func (iter *ServiceConfigIterator) Cursor() (datastore.Cursor, error) {
	return iter.iter.Cursor()
}

func (iter *ServiceConfigIterator) MustCursor() datastore.Cursor {
	c, err := iter.iter.Cursor()
	xerrors.MustNil(err)
	return c
}

func (iter *ServiceConfigIterator) Next() (*datastore.Key, *ServiceConfig, error) {
	if iter.viaKeys {
		key, err := iter.iter.Next(nil)
		if err != nil {
			if err == datastore.Done {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		_, ent, err := serviceConfigKindInstance.Get(iter.ctx, key)
		if err != nil {
			return nil, nil, err
		}
		return key, ent, nil

	}
	var ent ServiceConfig
	key, err := iter.iter.Next(&ent)
	if err != nil {
		if err == datastore.Done {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return key, &ent, nil
}

func (iter *ServiceConfigIterator) MustNext() (*datastore.Key, *ServiceConfig) {
	key, ent, err := iter.Next()
	xerrors.MustNil(err)
	return key, ent
}

var serviceConfigKindInstance = &ServiceConfigKind{}
