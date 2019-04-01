package asynctask

import (
	"context"

	"github.com/yssk22/go/keyvalue"
)

// Logic is an interface to execute a task
type Logic interface {
	Run(context.Context, *keyvalue.GetProxy, *AsyncTask) (*Progress, error)
}

// LogicFunc is an function to implement Logic
type LogicFunc func(context.Context, *keyvalue.GetProxy, *AsyncTask) (*Progress, error)

// Run implements Logic#Run
func (f LogicFunc) Run(ctx context.Context, params *keyvalue.GetProxy, t *AsyncTask) (*Progress, error) {
	return f(ctx, params, t)
}
