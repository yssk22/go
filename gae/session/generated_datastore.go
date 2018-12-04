// Code generated by github.com/yssk22/go/generator DO NOT EDIT.
//
package session

import (
	"context"
	ds "github.com/yssk22/go/gae/datastore"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xtime"
	"google.golang.org/appengine/datastore"
)

func (s *Session) NewKey(ctx context.Context) *datastore.Key {
	return ds.NewKey(ctx, "Session", s.ID)
}

type SessionReplacer interface {
	Replace(*Session, *Session) *Session
}

type SessionReplacerFunc func(*Session, *Session) *Session

func (f SessionReplacerFunc) Replace(old *Session, new *Session) *Session {
	return f(old, new)
}

type SessionKind struct{}

func NewSessionKind() *SessionKind {
	return sessionKindInstance
}

func (d *SessionKind) Get(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, *Session, error) {
	keys, ents, err := d.GetMulti(ctx, []interface{}{key}, options...)
	if err != nil {
		return nil, nil, err
	}
	return keys[0], ents[0], nil
}

func (d *SessionKind) MustGet(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, *Session) {
	k, v, e := d.Get(ctx, key, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *SessionKind) GetMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, []*Session, error) {
	var err error
	var dsKeys []*datastore.Key
	var ents []*Session
	if dsKeys, err = ds.NormalizeKeys(ctx, "Session", keys); err != nil {
		return nil, nil, xerrors.Wrap(err, "could not normalize keys: %v", keys)
	}
	size := len(dsKeys)
	if size == 0 {
		return nil, nil, nil
	}
	ents = make([]*Session, size, size)
	if err = ds.GetMulti(ctx, dsKeys, ents, options...); err != nil {
		return nil, nil, err
	}
	return dsKeys, ents, nil
}

func (d *SessionKind) MustGetMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, []*Session) {
	k, v, e := d.GetMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *SessionKind) Put(ctx context.Context, ent *Session, options ...ds.CRUDOption) (*datastore.Key, error) {
	keys, err := d.PutMulti(ctx, []*Session{ent}, options...)
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

func (d *SessionKind) MustPut(ctx context.Context, ent *Session, options ...ds.CRUDOption) *datastore.Key {
	k, e := d.Put(ctx, ent, options...)
	xerrors.MustNil(e)
	return k
}

func (d *SessionKind) PutMulti(ctx context.Context, ents []*Session, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	var err error
	var size = len(ents)
	var dsKeys []*datastore.Key
	dsKeys = make([]*datastore.Key, size, size)
	for i := range ents {
		dsKeys[i] = ents[i].NewKey(ctx)
		ents[i].Timestamp = xtime.Now()
	}
	if dsKeys, err = ds.PutMulti(ctx, dsKeys, ents); err != nil {
		return nil, err
	}
	return dsKeys, nil
}

func (d *SessionKind) MustPutMulti(ctx context.Context, ents []*Session, options ...ds.CRUDOption) []*datastore.Key {
	keys, err := d.PutMulti(ctx, ents, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *SessionKind) Delete(ctx context.Context, key interface{}, options ...ds.CRUDOption) (*datastore.Key, error) {
	keys, err := d.DeleteMulti(ctx, []interface{}{key}, options...)
	if err != nil {
		return nil, err
	}
	return keys[0], nil
}

func (d *SessionKind) MustDelete(ctx context.Context, key interface{}, options ...ds.CRUDOption) *datastore.Key {
	k, e := d.Delete(ctx, key, options...)
	xerrors.MustNil(e)
	return k
}

func (d *SessionKind) DeleteMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) ([]*datastore.Key, error) {
	var err error
	var dsKeys []*datastore.Key
	if dsKeys, err = ds.NormalizeKeys(ctx, "Session", keys); err != nil {
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

func (d *SessionKind) MustDeleteMulti(ctx context.Context, keys interface{}, options ...ds.CRUDOption) []*datastore.Key {
	k, e := d.DeleteMulti(ctx, keys, options...)
	xerrors.MustNil(e)
	return k
}

func (d *SessionKind) DeleteMatched(ctx context.Context, q *SessionQuery, options ...ds.CRUDOption) ([]*datastore.Key, error) {
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

func (d *SessionKind) MustDeleteMatched(ctx context.Context, q *SessionQuery, options ...ds.CRUDOption) []*datastore.Key {
	keys, err := d.DeleteMatched(ctx, q, options...)
	xerrors.MustNil(err)
	return keys
}

func (d *SessionKind) Replace(ctx context.Context, ent *Session, replacer SessionReplacer, options ...ds.CRUDOption) (*datastore.Key, *Session, error) {
	keys, ents, err := d.ReplaceMulti(ctx, []*Session{ent}, replacer, options...)
	if err != nil {
		return nil, ents[0], err
	}
	return keys[0], ents[0], err
}

func (d *SessionKind) MustReplace(ctx context.Context, ent *Session, replacer SessionReplacer, options ...ds.CRUDOption) (*datastore.Key, *Session) {
	k, v, e := d.Replace(ctx, ent, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

func (d *SessionKind) ReplaceMulti(ctx context.Context, ents []*Session, replacer SessionReplacer, options ...ds.CRUDOption) ([]*datastore.Key, []*Session, error) {
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

func (d *SessionKind) MustReplaceMulti(ctx context.Context, ents []*Session, replacer SessionReplacer, options ...ds.CRUDOption) ([]*datastore.Key, []*Session) {
	k, v, e := d.ReplaceMulti(ctx, ents, replacer, options...)
	xerrors.MustNil(e)
	return k, v
}

type SessionQuery struct {
	query   *ds.Query
	viaKeys bool
}

func NewSessionQuery() *SessionQuery {
	return &SessionQuery{
		query:   ds.NewQuery("Session"),
		viaKeys: false,
	}
}

func (d *SessionQuery) EqID(v string) *SessionQuery {
	d.query = d.query.Eq("ID", v)
	return d
}

func (d *SessionQuery) LtID(v string) *SessionQuery {
	d.query = d.query.Lt("ID", v)
	return d
}

func (d *SessionQuery) LeID(v string) *SessionQuery {
	d.query = d.query.Le("ID", v)
	return d
}

func (d *SessionQuery) GtID(v string) *SessionQuery {
	d.query = d.query.Gt("ID", v)
	return d
}

func (d *SessionQuery) GeID(v string) *SessionQuery {
	d.query = d.query.Ge("ID", v)
	return d
}

func (d *SessionQuery) NeID(v string) *SessionQuery {
	d.query = d.query.Ne("ID", v)
	return d
}

func (d *SessionQuery) AscID() *SessionQuery {
	d.query = d.query.Asc("ID")
	return d
}

func (d *SessionQuery) DescID() *SessionQuery {
	d.query = d.query.Desc("ID")
	return d
}

func (d *SessionQuery) Start(s string) *SessionQuery {
	d.query = d.query.Start(s)
	return d
}

func (d *SessionQuery) End(s string) *SessionQuery {
	d.query = d.query.End(s)
	return d
}

func (d *SessionQuery) Limit(n int) *SessionQuery {
	d.query = d.query.Limit(n)
	return d
}

func (d *SessionQuery) ViaKeys() *SessionQuery {
	d.viaKeys = true
	return d
}

func (d *SessionQuery) GetAll(ctx context.Context) ([]*datastore.Key, []Session, error) {
	if d.viaKeys {
		keys, err := d.query.KeysOnly().GetAll(ctx, nil)
		if err != nil {
			return nil, nil, err
		}
		_, ents, err := sessionKindInstance.GetMulti(ctx, keys)
		if err != nil {
			return nil, nil, err
		}
		list := make([]Session, len(ents))
		for i, e := range ents {
			list[i] = *e
		}
		return keys, list, nil
	}
	var ent []Session
	keys, err := d.query.GetAll(ctx, &ent)
	if err != nil {
		return nil, nil, err
	}
	return keys, ent, nil
}

func (d *SessionQuery) MustGetAll(ctx context.Context) ([]*datastore.Key, []Session) {
	keys, ents, err := d.GetAll(ctx)
	xerrors.MustNil(err)
	return keys, ents
}

func (d *SessionQuery) Count(ctx context.Context) (int, error) {
	return d.query.Count(ctx)
}

func (d *SessionQuery) MustCount(ctx context.Context) int {
	c, err := d.query.Count(ctx)
	xerrors.MustNil(err)
	return c
}

func (d *SessionQuery) Run(ctx context.Context) (*SessionIterator, error) {
	iter, err := d.query.Run(ctx)
	if err != nil {
		return nil, err
	}
	return &SessionIterator{
		ctx:     ctx,
		iter:    iter,
		viaKeys: d.viaKeys,
	}, err
}

func (d *SessionQuery) MustRun(ctx context.Context) *SessionIterator {
	iter, err := d.Run(ctx)
	xerrors.MustNil(err)
	return iter
}

func (d *SessionQuery) RunAll(ctx context.Context) ([]datastore.Key, []Session, string, error) {
	iter, err := d.Run(ctx)
	if err != nil {
		return nil, nil, "", err
	}
	var keys []datastore.Key
	var ents []Session
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

func (d *SessionQuery) MustRunAll(ctx context.Context) ([]datastore.Key, []Session, string) {
	keys, ents, next, err := d.RunAll(ctx)
	xerrors.MustNil(err)
	return keys, ents, next
}

type SessionIterator struct {
	ctx     context.Context
	iter    *datastore.Iterator
	viaKeys bool
}

func (iter *SessionIterator) Cursor() (datastore.Cursor, error) {
	return iter.iter.Cursor()
}

func (iter *SessionIterator) MustCursor() datastore.Cursor {
	c, err := iter.iter.Cursor()
	xerrors.MustNil(err)
	return c
}

func (iter *SessionIterator) Next() (*datastore.Key, *Session, error) {
	if iter.viaKeys {
		key, err := iter.iter.Next(nil)
		if err != nil {
			if err == datastore.Done {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		_, ent, err := sessionKindInstance.Get(iter.ctx, key)
		if err != nil {
			return nil, nil, err
		}
		return key, ent, nil

	}
	var ent Session
	key, err := iter.iter.Next(&ent)
	if err != nil {
		if err == datastore.Done {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return key, &ent, nil
}

func (iter *SessionIterator) MustNext() (*datastore.Key, *Session) {
	key, ent, err := iter.Next()
	xerrors.MustNil(err)
	return key, ent
}

var sessionKindInstance = &SessionKind{}
