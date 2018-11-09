package flow

import (
	"testing"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestFlow(t *testing.T) {
	a := assert.New(t)
	runner := generator.NewRunner(
		NewGenerator(),
	)
	a.Nil(runner.Run("./example"))
	// expected, _ := ioutil.ReadFile("./example/__generated.expected")
	// actual, _ := ioutil.ReadFile("./example/__generated__apis.go")
	// a.EqStr(string(expected), string(actual))
}
