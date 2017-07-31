package datastore

import (
	"fmt"
	"strings"

	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xlog"
	"context"
	"google.golang.org/appengine/datastore"
)

// Query is a struct to build a query to datastore.
type Query struct {
	kind        string
	ancestor    lazy.Value
	projects    []string
	orders      []*order
	filters     []*filter
	startCursor lazy.Value
	endCursor   lazy.Value
	limit       lazy.Value
	offset      lazy.Value
}

// NewQuery returns a new *Query for `kind`
func NewQuery(kind string) *Query {
	return &Query{
		kind: kind,
	}
}

// Ancestor sets the ancestor filter
func (q *Query) Ancestor(a lazy.Value) *Query {
	q.ancestor = a
	return q
}

// Eq sets the "=" filter on the `name` field.
func (q *Query) Eq(name string, value lazy.Value) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeEq,
		Value: value,
	})
	return q
}

// Lt sets the `<` filter on the `name` field.
func (q *Query) Lt(name string, value lazy.Value) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeLt,
		Value: value,
	})
	return q
}

// Le sets the `<=` filter on the `name` field.
func (q *Query) Le(name string, value lazy.Value) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeLe,
		Value: value,
	})
	return q
}

// Gt sets the `>` filter on the `name` field.
func (q *Query) Gt(name string, value lazy.Value) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeGt,
		Value: value,
	})
	return q
}

// Ge sets the `>=` filter on the `name` field.
func (q *Query) Ge(name string, value lazy.Value) *Query {
	q.filters = append(q.filters, &filter{
		Name:  name,
		Type:  filterTypeGe,
		Value: value,
	})
	return q
}

// Ne sets the `!=` filter on the `name` field.
func (q *Query) Ne(name string, value lazy.Value) *Query {
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
func (q *Query) Limit(value lazy.Value) *Query {
	q.limit = value
	return q
}

// Start sets the start cursor
func (q *Query) Start(value lazy.Value) *Query {
	q.startCursor = value
	return q
}

// End sets the end cursor
func (q *Query) End(value lazy.Value) *Query {
	q.endCursor = value
	return q
}

type order struct {
	Name string
	Type orderType
}

type orderType string

const (
	orderTypeAsc  orderType = ""
	orderTypeDesc           = "-"
)

type filter struct {
	Name  string
	Type  filterType
	Value lazy.Value
}

type filterType string

const (
	filterTypeEq filterType = "="
	filterTypeLe            = "<="
	filterTypeLt            = "<"
	filterTypeGe            = ">="
	filterTypeGt            = ">"
	filterTypeNe            = "!="
)

// FilterValueOmit is a value to omit the filter.
//
//     q.Eq("FieladName", lazy.LazyFunc(func(ctx) lazy.Value {
//         v, ok :=  ctx.Value("filterValue").(string)
//         if ok {
//             return v;
//         }
//         return query.FilterValueOmit
//     }))
//
// This means if 'filterValue' does not present in the context nor a int value,
// FieldName filter is not be applied in the query q.
var FilterValueOmit = lazy.Value(nil)

// GetAll fills the query result into dst and returns corresponding *datastore.Key
func (q *Query) GetAll(ctx context.Context, dst interface{}) ([]*datastore.Key, error) {
	query, err := q.prepare(ctx)
	if err != nil {
		return nil, err
	}
	return query.GetAll(ctx, dst)
}

// Run runs a query and returns *datastore.Iterator
func (q *Query) Run(ctx context.Context) (*datastore.Iterator, error) {
	query, err := q.prepare(ctx)
	if err != nil {
		return nil, err
	}
	return query.Run(ctx), nil
}

// Count returns a count
func (q *Query) Count(ctx context.Context) (int, error) {
	query, err := q.prepare(ctx)
	if err != nil {
		return 0, err
	}
	return query.Count(ctx)
}

func (q *Query) prepare(ctx context.Context) (*datastore.Query, error) {
	var buff []string
	logger := xlog.WithContext(ctx).WithKey(LoggerKey)
	query := datastore.NewQuery(q.kind)

	if q.projects != nil {
		query = query.Project(q.projects...)
		buff = append(buff, fmt.Sprintf("Projects: %v", q.projects))
	}

	if q.ancestor != nil {
		v, err := q.ancestor.Eval(ctx)
		if err != nil {
			return nil, fmt.Errorf("ancestor field error: %v", err)
		}
		if vv, ok := v.(*datastore.Key); ok {
			query = query.Ancestor(vv)
			buff = append(buff, fmt.Sprintf("Ancestor: %v", vv))
		}
	}

	if q.filters != nil {
		var s []string
		for _, f := range q.filters {
			val, err := f.Value.Eval(ctx)
			if err != nil {
				return nil, fmt.Errorf("%s filter error: %v", f.Name, err)
			}
			if val != FilterValueOmit {
				query = query.Filter(
					fmt.Sprintf("%s %s", f.Name, f.Type),
					val,
				)
				s = append(s, fmt.Sprintf("%s %s %s", f.Name, f.Type, val))
			}
		}
		buff = append(buff, fmt.Sprintf("Filter: %v", strings.Join(s, " AND ")))
	}

	if q.orders != nil {
		var s []string
		for _, o := range q.orders {
			order := fmt.Sprintf("%s%s", o.Type, o.Name)
			query = query.Order(
				fmt.Sprintf(order),
			)
			s = append(s, order)
		}
		buff = append(buff, fmt.Sprintf("Order: %v", strings.Join(s, " ")))
	}

	if q.limit != nil {
		v, err := q.limit.Eval(ctx)
		if err != nil {
			return nil, fmt.Errorf("limit error: %v", err)
		}
		if vv, ok := v.(int); ok {
			query = query.Limit(vv)
			buff = append(buff, fmt.Sprintf("Limit: %d", vv))
		}
	}

	if q.startCursor != nil {
		v, err := q.startCursor.Eval(ctx)
		if err != nil {
			return nil, fmt.Errorf("start cursor error: %v", err)
		}
		if str := fmt.Sprintf("%s", v); str != "" {
			if vv, err := datastore.DecodeCursor(str); err == nil {
				query = query.Start(vv)
				buff = append(buff, fmt.Sprintf("Start: %s", vv))
			}
		}
	}

	if q.endCursor != nil {
		v, err := q.endCursor.Eval(ctx)
		if err != nil {
			return nil, fmt.Errorf("end cursor error: %v", err)
		}
		if str := fmt.Sprintf("%s", v); str != "" {
			if vv, err := datastore.DecodeCursor(str); err == nil {
				query = query.End(vv)
				buff = append(buff, fmt.Sprintf("End: %s", vv))
			}
		}
	}
	logger.Debug(func(p *xlog.Printer) {
		p.Printf("Query: Kind=%s\n", q.kind)
		for _, line := range buff {
			p.Printf("\t%s\n", line)
		}
	})
	return query, nil
}
