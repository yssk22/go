package httptest

import (
	"testing"

	"github.com/yssk22/go/x/xtesting"
	"github.com/yssk22/go/x/xtesting/assert"
)

// Runner is a struct to run a test
type Runner struct {
	r *xtesting.Runner
}

// NewRunner returns a *Runner
func NewRunner(t *testing.T) *Runner {
	return &Runner{
		r: xtesting.NewRunner(t),
	}
}

// Setup wraps xtesting.Runner#Setup
func (r *Runner) Setup(f func(a *Assert)) {
	r.r.Setup(func(a *assert.Assert) {
		f(&Assert{a})
	})
}

// Teardown wraps xtesting.Runner#Teardown
func (r *Runner) Teardown(f func(a *Assert)) {
	r.r.Teardown(func(a *assert.Assert) {
		f(&Assert{a})
	})
}

// Run wraps xtesting.Runner#Run
func (r *Runner) Run(name string, f func(a *Assert)) {
	r.r.Run(name, func(a *assert.Assert) {
		f(&Assert{a})
	})
}
