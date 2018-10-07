package example

import (
	"github.com/yssk22/go/tools/cmd/ent/example"
	"github.com/yssk22/go/web"
)

// Example struct
// @jstype
type Example struct {
	ID string `json:"id"`
}

// @api path=/path/to/example
func getExample(req *web.Request) (*Example, error) {
	a := &Example{
		ID: "myid",
	}
	return a, nil
}

// @api path=/path/to/example2
func getExample2(req *web.Request) (*example.Example, error) {
	a := &example.Example{
		ID: "myid",
	}
	return a, nil
}
