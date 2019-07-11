package datastore

import (
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
)

type orderType string

const (
	orderTypeAsc  orderType = ""
	orderTypeDesc           = "-"
)

// Query is a wrapper for datasatore.Query
type Query struct {
	inner     *datastore.Query
	kind      string
	statement []string // for debugging
}

// NewQuery returns a *Query for the kind k
func NewQuery(k string) *Query {
	return &Query{
		inner:     datastore.NewQuery(k),
		kind:      k,
		statement: []string{},
	}
}

// Namespace sets the namespace filter
func (q *Query) Namespace(ns string) *Query {
	q.inner = q.inner.Namespace(ns)
	q.statement = append(q.statement, fmt.Sprintf("Namespace(%s)", ns))
	return q
}

// Ancestor sets the ancestor filter
func (q *Query) Ancestor(a *datastore.Key) *Query {
	q.inner = q.inner.Ancestor(a)
	q.statement = append(q.statement, fmt.Sprintf("Ancestor(%s)", a))
	return q
}

// KeysOnly sets the query to return only keys
func (q *Query) KeysOnly() *Query {
	q.inner = q.inner.KeysOnly()
	q.statement = append(q.statement, "KeysOnly")
	return q
}

// Eq sets the "=" filter on the `name` field.
func (q *Query) Eq(name string, value interface{}) *Query {
	q.inner = q.inner.Filter(fmt.Sprintf("%s =", name), value)
	q.statement = append(q.statement, fmt.Sprintf("%s = %s", name, value))
	return q
}

// Lt sets the `<` filter on the `name` field.
func (q *Query) Lt(name string, value interface{}) *Query {
	q.inner = q.inner.Filter(fmt.Sprintf("%s <", name), value)
	q.statement = append(q.statement, fmt.Sprintf("%s < %s", name, value))
	return q
}

// Le sets the `<=` filter on the `name` field.
func (q *Query) Le(name string, value interface{}) *Query {
	q.inner = q.inner.Filter(fmt.Sprintf("%s <=", name), value)
	q.statement = append(q.statement, fmt.Sprintf("%s <= %s", name, value))
	return q
}

// Gt sets the `>` filter on the `name` field.
func (q *Query) Gt(name string, value interface{}) *Query {
	q.inner = q.inner.Filter(fmt.Sprintf("%s >", name), value)
	q.statement = append(q.statement, fmt.Sprintf("%s > %s", name, value))
	return q
}

// Ge sets the `>=` filter on the `name` field.
func (q *Query) Ge(name string, value interface{}) *Query {
	q.inner = q.inner.Filter(fmt.Sprintf("%s >=", name), value)
	q.statement = append(q.statement, fmt.Sprintf("%s >= %s", name, value))
	return q
}

// Ne sets the `!=` filter on the `name` field.
func (q *Query) Ne(name string, value interface{}) *Query {
	q.inner = q.inner.Filter(fmt.Sprintf("%s !=", name), value)
	q.statement = append(q.statement, fmt.Sprintf("%s != %s", name, value))
	return q
}

// Asc specifies ascending order on the given filed.
func (q *Query) Asc(name string) *Query {
	q.inner = q.inner.Order(name)
	q.statement = append(q.statement, fmt.Sprintf("Asc(%s)", name))
	return q
}

// Desc specifies descending order on the given filed.
func (q *Query) Desc(name string) *Query {
	q.inner = q.inner.Order(fmt.Sprintf("-%s", name))
	q.statement = append(q.statement, fmt.Sprintf("Desc(%s)", name))
	return q
}

// Limit sets the limit
func (q *Query) Limit(value int) *Query {
	q.inner = q.inner.Limit(value)
	q.statement = append(q.statement, fmt.Sprintf("Limit(%d)", value))
	return q
}

// Start sets the start cursor
func (q *Query) Start(value string) *Query {
	if c, err := datastore.DecodeCursor(value); err == nil {
		q.inner = q.inner.Start(c)
		q.statement = append(q.statement, fmt.Sprintf("Start(%q)", value))
	}
	return q
}

// End sets the end cursor
func (q *Query) End(value string) *Query {
	if c, err := datastore.DecodeCursor(value); err == nil {
		q.inner = q.inner.End(c)
		q.statement = append(q.statement, fmt.Sprintf("End(%q)", value))
	}
	return q
}

func (q *Query) String() string {
	return fmt.Sprintf("Query[%s](%s)", q.kind, strings.Join(q.statement, ", "))
}
