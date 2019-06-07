package datastore

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/yssk22/go/x/xlog"
)

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

// Query is a wrapper for datasatore.Query
type Query struct {
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

// NewQuery returns a *Query for the kind k
func NewQuery(k string) *Query {
	return &Query{
		kind: k,
	}
}

// Ancestor sets the ancestor filter
func (q *Query) Ancestor(a *datastore.Key) *Query {
	q.ancestor = a
	return q
}

// KeysOnly sets the query to return only keys
func (q *Query) KeysOnly() *Query {
	q.keysOnly = true
	return q
}

// Eq sets the "=" filter on the `name` field.
func (q *Query) Eq(name string, value interface{}) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeEq,
		Value: value,
	})
	return q
}

// Lt sets the `<` filter on the `name` field.
func (q *Query) Lt(name string, value interface{}) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeLt,
		Value: value,
	})
	return q
}

// Le sets the `<=` filter on the `name` field.
func (q *Query) Le(name string, value interface{}) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeLe,
		Value: value,
	})
	return q
}

// Gt sets the `>` filter on the `name` field.
func (q *Query) Gt(name string, value interface{}) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeGt,
		Value: value,
	})
	return q
}

// Ge sets the `>=` filter on the `name` field.
func (q *Query) Ge(name string, value interface{}) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeGe,
		Value: value,
	})
	return q
}

// Ne sets the `!=` filter on the `name` field.
func (q *Query) Ne(name string, value interface{}) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeNe,
		Value: value,
	})
	return q
}

// Asc specifies ascending order on the given filed.
func (q *Query) Asc(name string) *Query {
	q.orders = append(q.orders, &order{
		Name: name,
		Type: orderTypeAsc,
	})
	return q
}

// Desc specifies descending order on the given filed.
func (q *Query) Desc(name string) *Query {
	q.orders = append(q.orders, &order{
		Name: name,
		Type: orderTypeDesc,
	})
	return q
}

// Limit sets the limit
func (q *Query) Limit(value int) *Query {
	q.limit = value
	return q
}

// Start sets the start cursor
func (q *Query) Start(value string) *Query {
	q.startCursor = value
	return q
}

// End sets the end cursor
func (q *Query) End(value string) *Query {
	q.endCursor = value
	return q
}

// WithLogging sets the logger key for the query
func (q *Query) WithLogging(loggerKey string) *Query {
	q.loggerKey = loggerKey
	return q
}

// ViaKeys runs GetAll query to get keys only at first then call GetMulti with these keys to utilize the cache.ViaKeys
func (q *Query) ViaKeys() *Query {
	q.viaKeys = true
	return q
}

func (q *Query) prepare(ctx context.Context) (*datastore.Query, error) {
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
