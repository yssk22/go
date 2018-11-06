package validator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_Required(t *testing.T) {
	a := assert.New(t)
	type T struct{}
	var tt *T
	a.Nil(Required(1))
	a.Nil(Required(struct{}{}))
	a.NotNil(Required(nil))
	a.NotNil(Required(tt))
}
