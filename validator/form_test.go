package validator

import (
	"net/url"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestFormValidation(t *testing.T) {
	a := assert.New(t)
	obj := url.Values{
		"foo": []string{"bar"},
	}

	v := NewFormValidator()
	v.Field("foo").Required()
	v.Field("var")
	a.Nil(v.Eval(obj))

	obj["foo"] = []string{}
	a.NotNil(v.Eval(obj))

	obj["foo"] = nil
	a.NotNil(v.Eval(obj))

	delete(obj, "foo")
	a.NotNil(v.Eval(obj))
}

func TestFormValidation_IntField(t *testing.T) {
	a := assert.New(t)
	obj := url.Values{
		"foo": []string{"0"},
	}

	v := NewFormValidator()
	v.IntField("foo").Required().Max(10)
	a.Nil(v.Eval(obj))

	obj["foo"] = []string{"100"}
	a.NotNil(v.Eval(obj))

	obj["foo"] = []string{}
	a.NotNil(v.Eval(obj))

	obj["foo"] = nil
	a.NotNil(v.Eval(obj))

	delete(obj, "foo")
	a.NotNil(v.Eval(obj))

}
