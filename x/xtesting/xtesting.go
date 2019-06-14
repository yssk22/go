package xtesting

import (
	"flag"
	"strings"
	"testing"

	"github.com/yssk22/go/x/xruntime"
	"github.com/yssk22/go/x/xtesting/assert"
)

var isTesting = false

func init() {
	isTesting = (flag.Lookup("test.v") != nil)
}

// IsTesting returns if the current process is executed by `go test` or not.
func IsTesting() bool {
	return isTesting
}

// Runner is a struct to run a test
type Runner struct {
	t        *testing.T
	setup    func(a *assert.Assert)
	teardown func(a *assert.Assert)
}

// NewRunner returns a *Runner
func NewRunner(t *testing.T) *Runner {
	return &Runner{
		t: t,
	}
}

// Setup sets the function to be executed on every test execution
func (r *Runner) Setup(f func(a *assert.Assert)) {
	r.setup = f
}

// Teardown sets the function to be executed on every test execution
func (r *Runner) Teardown(f func(a *assert.Assert)) {
	r.teardown = f
}

// Run runs a test
func (r *Runner) Run(name string, f func(a *assert.Assert)) {
	r.t.Run(name, func(t *testing.T) {
		a := assert.New(t)
		defer func() {
			if err := recover(); err != nil {
				if r.teardown != nil {
					defer func() {
						if err := recover(); err != nil {
							stacks := xruntime.CollectAllStacksSimple()
							t.Errorf("panic: %v\n%s", err, strings.Join(stacks, "\n"))
						}
					}()
					r.teardown(a)
					stacks := xruntime.CollectAllStacksSimple()
					t.Errorf("panic: %v\n%s", err, strings.Join(stacks, "\n"))
				} else {
					stacks := xruntime.CollectAllStacksSimple()
					t.Errorf("panic: %v\n%s", err, strings.Join(stacks, "\n"))
				}
			}
		}()
		if r.setup != nil {
			r.setup(a)
		}
		f(a)
	})
}
