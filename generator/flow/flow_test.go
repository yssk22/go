package flow

import (
	"io/ioutil"
	"testing"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestFlow(t *testing.T) {
	a := assert.New(t)
	runner := generator.NewRunner(
		NewGenerator(NewOptions()),
	)
	a.Nil(runner.Run("./example"))
	expected, _ := ioutil.ReadFile("./example/GoTypes.expected.js")
	actual, _ := ioutil.ReadFile("./example/GoTypes.js")
	a.EqStr(string(expected), string(actual))
}
