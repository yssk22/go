package datastore

import (
	"context"
	"fmt"
	"strings"

	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xlog"
	"google.golang.org/appengine/datastore"
)

// QueryOption is an option arg to build a query
type QueryOption func(*query) *query

// NewQuery returns an *query configured with options
func NewQuery(kind string) Query {
	return &query{
		kind: kind,
	}
}

// Query is an interface to run a query
type Query interface {
	// execution functions

	GetAll(context.Context, interface{}) ([]*datastore.Key, error)
	MustGetAll(context.Context, interface{}) []*datastore.Key
	Run(context.Context) (*datastore.Iterator, error)
	MustRun(context.Context) *datastore.Iterator
	Count(context.Context) (int, error)
	MustCount(context.Context) int

	// builder functions
	Ancestor(*datastore.Key) Query
	KeysOnly(bool) Query
	Eq(string, interface{}) Query
	Lt(string, interface{}) Query
	Le(string, interface{}) Query
	Gt(string, interface{}) Query
	Ge(string, interface{}) Query
	Ne(string, interface{}) Query
	Asc(string) Query
	Desc(string) Query
	Limit(int) Query
	Start(string) Query
	End(string) Query
	WithLogging(key string) Query
}

type orderType string

const (
	orderTypeAsc  orderType = ""
	orderTypeDesc           = "-"
)

type filterType string

const (
	filterTypeEq filterType = "="
	filterTypeLe            = "<="
	filterTypeLt            = "<"
	filterTypeGe            = ">="
	filterTypeGt            = ">"
	filterTypeNe            = "!="
)

type order struct {
	Name string
	Type orderType
}

type filter struct {
	Name  string
	Type  filterType
	Value interface{}
}

type query struct {
	kind        string
	ancestor    *datastore.Key
	projects    []string
	keysOnly    bool
	orders      []*order
	filters     []*filter
	startCursor string
	endCursor   string
	limit       int
	offset      int
	loggerKey   string `json:"-"`
}

// Ancestor sets the ancestor filter
func (q *query) Ancestor(a *datastore.Key) Query {
	q.ancestor = a
	return q
}

// KeysOnly sets the query to fetch only keys.
func (q *query) KeysOnly(t bool) Query {
	q.keysOnly = t
	return q
}

// Eq sets the "=" filter on the `name` field.
func (q *query) Eq(name string, value interface{}) Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeEq,
		Value: value,
	})
	return q
}

// Lt sets the `<` filter on the `name` field.
func (q *query) Lt(name string, value interface{}) Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeLt,
		Value: value,
	})
	return q
}

// Le sets the `<=` filter on the `name` field.
func (q *query) Le(name string, value interface{}) Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeLe,
		Value: value,
	})
	return q
}

// Gt sets the `>` filter on the `name` field.
func (q *query) Gt(name string, value interface{}) Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeGt,
		Value: value,
	})
	return q
}

// Ge sets the `>=` filter on the `name` field.
func (q *query) Ge(name string, value interface{}) Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeGe,
		Value: value,
	})
	return q
}

// Ne sets the `!=` filter on the `name` field.
func (q *query) Ne(name string, value interface{}) Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeNe,
		Value: value,
	})
	return q
}

// Asc specifies ascending order on the given filed.
func (q *query) Asc(name string) Query {
	q.orders = append(q.orders, &order{
		Name: name,
		Type: orderTypeAsc,
	})
	return q
}

// Desc specifies descending order on the given filed.
func (q *query) Desc(name string) Query {
	q.orders = append(q.orders, &order{
		Name: name,
		Type: orderTypeDesc,
	})
	return q
}

// Limit sets the limit
func (q *query) Limit(value int) Query {
	q.limit = value
	return q
}

// Start sets the start cursor
func (q *query) Start(value string) Query {
	q.startCursor = value
	return q
}

// End sets the end cursor
func (q *query) End(value string) Query {
	q.endCursor = value
	return q
}

// WithLogging sets the logger key for the query
func (q *query) WithLogging(loggerKey string) Query {
	q.loggerKey = loggerKey
	return q
}

// GetAll fills the query result into dst and returns corresponding *datastore.Key
func (q *query) GetAll(ctx context.Context, dst interface{}) ([]*datastore.Key, error) {
	prepared, err := q.prepare(ctx)
	if err != nil {
		return nil, err
	}
	return prepared.GetAll(ctx, dst)
}

// MustGetAll is like GetAll but panic if an error occurs
func (q *query) MustGetAll(ctx context.Context, dst interface{}) []*datastore.Key {
	keys, err := q.GetAll(ctx, dst)
	xerrors.MustNil(err)
	return keys
}

// Run runs a query and returns *datastore.Iterator
func (q *query) Run(ctx context.Context) (*datastore.Iterator, error) {
	prepared, err := q.prepare(ctx)
	if err != nil {
		return nil, err
	}
	return prepared.Run(ctx), nil
}

// MustRun is like Run but panic if an error occurs
func (q *query) MustRun(ctx context.Context) *datastore.Iterator {
	iter, err := q.Run(ctx)
	xerrors.MustNil(err)
	return iter
}

// Count returns a count
func (q *query) Count(ctx context.Context) (int, error) {
	prepared, err := q.prepare(ctx)
	if err != nil {
		return 0, err
	}
	return prepared.Count(ctx)
}

// MustCount is like Count but panic if an error occurs
func (q *query) MustCount(ctx context.Context) int {
	c, err := q.Count(ctx)
	xerrors.MustNil(err)
	return c
}

func (q *query) prepare(ctx context.Context) (*datastore.Query, error) {
	var logger *xlog.Logger
	var buff []string
	if q.loggerKey != "" {
		defer func() {
			ctx, logger = xlog.WithContextAndKey(ctx, "", q.loggerKey)
			logger.Info(func(p *xlog.Printer) {
				p.Printf("Query: Kind=%s\n", q.kind)
				for _, line := range buff {
					p.Printf("\t%s\n", line)
				}
			})
		}()
	}

	query := datastore.NewQuery(q.kind)
	if q.projects != nil {
		query = query.Project(q.projects...)
		buff = append(buff, fmt.Sprintf("Projects: %v", q.projects))
	}

	if q.ancestor != nil {
		query = query.Ancestor(q.ancestor)
		buff = append(buff, fmt.Sprintf("Ancestor: %v", q.ancestor))
	}

	if q.filters != nil {
		var s []string
		for _, f := range q.filters {
			query = query.Filter(
				fmt.Sprintf("%s %s", f.Name, f.Type),
				f.Value,
			)
			s = append(s, fmt.Sprintf("%s %s %s", f.Name, f.Type, f.Value))
		}
		buff = append(buff, fmt.Sprintf("Filter: %v", strings.Join(s, " AND ")))
	}

	if q.orders != nil {
		var s []string
		for _, o := range q.orders {
			order := fmt.Sprintf("%s%s", o.Type, o.Name)
			query = query.Order(order)
			s = append(s, order)
		}
		buff = append(buff, fmt.Sprintf("Order: %v", strings.Join(s, " ")))
	}

	if q.limit > 0 {
		query = query.Limit(q.limit)
		buff = append(buff, fmt.Sprintf("Limit: %d", q.limit))
	}

	if q.startCursor != "" {
		if vv, err := datastore.DecodeCursor(q.startCursor); err == nil {
			query = query.Start(vv)
			buff = append(buff, fmt.Sprintf("Start: %s", vv))
		}
	}

	if q.endCursor != "" {
		if vv, err := datastore.DecodeCursor(q.endCursor); err == nil {
			query = query.End(vv)
			buff = append(buff, fmt.Sprintf("End: %s", vv))
		}
	}

	if q.keysOnly {
		query = query.KeysOnly()
		buff = append(buff, "KeysOnly: true")
	}
	return query, nil
}
