package enum

import (
	"io/ioutil"
	"testing"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestEnum(t *testing.T) {
	a := assert.New(t)
	runner := generator.NewRunner(
		NewGenerator(),
	)
	a.Nil(runner.Run("./example"))
	expected, _ := ioutil.ReadFile("./example/__generated.expected")
	actual, _ := ioutil.ReadFile("./example/generated_enums.go")
	a.EqStr(string(expected), string(actual))
}
