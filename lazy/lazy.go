// Package lazy provides context-based lazy evaluatios
package lazy

import (
	"golang.org/x/net/context"
)

// Value is an interface to represent lazy value
type Value interface {
	Eval(context.Context) (interface{}, error)
}

// Func is an alias to lazy evaluated function
type Func func(context.Context) (interface{}, error)

// Eval implements Value#Eval
func (f Func) Eval(ctx context.Context) (interface{}, error) {
	return f(ctx)
}

// New returns a new Value to return v
func New(v interface{}) Value {
	return Func(func(context.Context) (interface{}, error) {
		return v, nil
	})
}
