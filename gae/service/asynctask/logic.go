package asynctask

import (
	"github.com/speedland/go/keyvalue"
	"golang.org/x/net/context"
)

// Logic is an interface to execute a task
type Logic interface {
	Run(context.Context, *keyvalue.GetProxy, *AsyncTask) (*Progress, error)
}

// Func is an function to implement Logic
type Func func(context.Context, *keyvalue.GetProxy, *AsyncTask) (*Progress, error)

// Run implements Logic#Run
func (f Func) Run(ctx context.Context, params *keyvalue.GetProxy, t *AsyncTask) (*Progress, error) {
	return f(ctx, params, t)
}
