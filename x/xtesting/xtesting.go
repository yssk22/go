package xtesting

import (
	"flag"
	"testing"

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
							t.Errorf("panic detected on Teardown: %v", err)
						}
					}()
					t.Logf("%s:Teardown", t.Name())
					r.teardown(a)
					t.Errorf("panic detected on Setup or test itself: %v", err)
				}
			}
		}()
		if r.setup != nil {
			t.Logf("%s:Setup", t.Name())
			r.setup(a)
		}
		f(a)
	})
}
