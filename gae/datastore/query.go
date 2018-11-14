package datastore

import (
	"context"
	"fmt"
	"strings"

	"github.com/yssk22/go/x/xlog"
	"google.golang.org/appengine/datastore"
)

// QueryOption is a function to configure query operation
type QueryOption func(*query) *query

// Ancestor sets the ancestor filter
func Ancestor(a *datastore.Key) QueryOption {
	return QueryOption(func(q *query) *query {
		q.ancestor = a
		return q
	})
}

// Eq sets the "=" filter on the `name` field.
func Eq(name string, value interface{}) QueryOption {
	return QueryOption(func(q *query) *query {
		q.filters = append(q.filters, &filter{
			Name:  name,
			Type:  filterTypeEq,
			Value: value,
		})
		return q
	})
}

// Lt sets the `<` filter on the `name` field.
func Lt(name string, value interface{}) QueryOption {
	return QueryOption(func(q *query) *query {
		q.filters = append(q.filters, &filter{
			Name:  name,
			Type:  filterTypeLt,
			Value: value,
		})
		return q
	})
}

// Le sets the `<=` filter on the `name` field.
func Le(name string, value interface{}) QueryOption {
	return QueryOption(func(q *query) *query {
		q.filters = append(q.filters, &filter{
			Name:  name,
			Type:  filterTypeLe,
			Value: value,
		})
		return q
	})
}

// Gt sets the `>` filter on the `name` field.
func Gt(name string, value interface{}) QueryOption {
	return QueryOption(func(q *query) *query {
		q.filters = append(q.filters, &filter{
			Name:  name,
			Type:  filterTypeGt,
			Value: value,
		})
		return q
	})
}

// Ge sets the `>=` filter on the `name` field.
func Ge(name string, value interface{}) QueryOption {
	return QueryOption(func(q *query) *query {
		q.filters = append(q.filters, &filter{
			Name:  name,
			Type:  filterTypeGe,
			Value: value,
		})
		return q
	})
}

// Ne sets the `!=` filter on the `name` field.
func Ne(name string, value interface{}) QueryOption {
	return QueryOption(func(q *query) *query {
		q.filters = append(q.filters, &filter{
			Name:  name,
			Type:  filterTypeNe,
			Value: value,
		})
		return q
	})
}

// Asc specifies ascending order on the given filed.
func Asc(name string) QueryOption {
	return QueryOption(func(q *query) *query {
		q.orders = append(q.orders, &order{
			Name: name,
			Type: orderTypeAsc,
		})
		return q
	})
}

// Desc specifies descending order on the given filed.
func Desc(name string) QueryOption {
	return QueryOption(func(q *query) *query {
		q.orders = append(q.orders, &order{
			Name: name,
			Type: orderTypeDesc,
		})
		return q
	})
}

// Limit sets the limit
func Limit(value int) QueryOption {
	return QueryOption(func(q *query) *query {
		q.limit = value
		return q
	})
}

// Start sets the start cursor
func Start(value string) QueryOption {
	return QueryOption(func(q *query) *query {
		q.startCursor = value
		return q
	})
}

// End sets the end cursor
func End(value string) QueryOption {
	return QueryOption(func(q *query) *query {
		q.endCursor = value
		return q
	})
}

// WithLogging sets the logger key for the query
func WithLogging(loggerKey string) QueryOption {
	return QueryOption(func(q *query) *query {
		q.loggerKey = loggerKey
		return q
	})
}

// ViaKeys runs GetAll query to get keys only at first then call GetMulti with these keys to utilize the cache.ViaKeys
func ViaKeys() QueryOption {
	return QueryOption(func(q *query) *query {
		q.viaKeys = true
		return q
	})
}

// GetAll fills the query result into dst and returns corresponding *datastore.Key
func GetAll(ctx context.Context, kind string, dst interface{}, options ...QueryOption) ([]*datastore.Key, error) {
	q := &query{
		kind: kind,
	}
	for _, f := range options {
		q = f(q)
	}
	if dst == nil {
		q.keysOnly = true
	}
	if q.viaKeys {
		q.keysOnly = true
	}
	prepared, err := q.prepare(ctx)
	if err != nil {
		return nil, err
	}
	keys, err := prepared.GetAll(ctx, dst)
	if err != nil {
		return nil, err
	}
	if dst != nil && q.viaKeys {
		err = GetMulti(ctx, keys, dst)
		return keys, err
	}
	return keys, err
}

// Run runs a query and returns *datastore.Iterator
func Run(ctx context.Context, kind string, options ...QueryOption) (*datastore.Iterator, error) {
	q := &query{
		kind: kind,
	}
	for _, f := range options {
		q = f(q)
	}
	prepared, err := q.prepare(ctx)
	if err != nil {
		return nil, err
	}
	return prepared.Run(ctx), nil
}

// Count returns a count
func Count(ctx context.Context, kind string, options ...QueryOption) (int, error) {
	q := &query{
		kind: kind,
	}
	for _, f := range options {
		q = f(q)
	}
	prepared, err := q.prepare(ctx)
	if err != nil {
		return 0, err
	}
	return prepared.Count(ctx)
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
	loggerKey   string
	viaKeys     bool
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
