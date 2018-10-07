package api

import (
	"io/ioutil"
	"testing"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestAPI(t *testing.T) {
	a := assert.New(t)
	runner := generator.NewRunner(
		NewGenerator(),
	)
	runner.Run("./example")
	expected, _ := ioutil.ReadFile("./example/api.go.expected")
	actual, _ := ioutil.ReadFile("./example/api.go")
	a.EqStr(string(expected), string(actual))

}
