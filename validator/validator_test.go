package validator

import (
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestRequired(t *testing.T) {
	a := assert.New(t)
	type TestObject struct {
		Str string
	}

	obj := &TestObject{}
	v := NewValidator()
	v.Field("Str").Required()
	result := v.Eval(obj)
	a.NotNil(result)
	a.NotNil(result.Errors["Str"])
	a.EqStr("must be required", result.Errors["Str"][0].String())

	obj.Str = "exists"
	a.Nil(v.Eval(obj))
}

func TestMin(t *testing.T) {
	a := assert.New(t)
	type TestObject struct {
		Str   string
		Int   int
		Float float32
		Array []string
	}
	obj := &TestObject{}
	v := NewValidator()
	v.Field("Str").Min(1)
	v.Field("Int").Min(1)
	v.Field("Float").Min(1)
	v.Field("Array").Min(1)
	result := v.Eval(obj)
	a.NotNil(result.Errors["Str"])
	a.NotNil(result.Errors["Int"])
	a.NotNil(result.Errors["Float"])
	a.NotNil(result.Errors["Array"])
	a.EqStr("must be more than or equal to 1", result.Errors["Str"][0].String())
	a.EqStr("must be more than or equal to 1", result.Errors["Int"][0].String())
	a.EqStr("must be more than or equal to 1", result.Errors["Float"][0].String())
	a.EqStr("must be more than or equal to 1", result.Errors["Array"][0].String())

	obj.Str = "Foo"
	obj.Int = 5
	obj.Float = 5.5
	obj.Array = []string{"a", "b"}
	a.Nil(v.Eval(obj))
}

func TestMax(t *testing.T) {
	a := assert.New(t)
	type TestObject struct {
		Str   string
		Int   int
		Float float32
		Array []string
	}
	obj := &TestObject{}
	v := NewValidator()
	v.Field("Str").Max(1)
	v.Field("Int").Max(1)
	v.Field("Float").Max(1)
	v.Field("Array").Max(1)

	obj.Str = "Foo"
	obj.Int = 5
	obj.Float = 5
	obj.Array = []string{"a", "b", "c", "d"}
	result := v.Eval(obj)
	a.NotNil(result.Errors["Str"])
	a.NotNil(result.Errors["Int"])
	a.NotNil(result.Errors["Float"])
	a.NotNil(result.Errors["Array"])
	a.EqStr("must be less than or equal to 1", result.Errors["Str"][0].String())
	a.EqStr("must be less than or equal to 1", result.Errors["Int"][0].String())
	a.EqStr("must be less than or equal to 1", result.Errors["Float"][0].String())
	a.EqStr("must be less than or equal to 1", result.Errors["Array"][0].String())

	obj.Str = "F"
	obj.Int = 1
	obj.Float = 1
	obj.Array = []string{"a"}
	a.Nil(v.Eval(obj))
}

func TestMatch(t *testing.T) {
	a := assert.New(t)
	type TestObject struct {
		Str   string
		Bytes []byte
	}
	obj := &TestObject{}
	v := NewValidator()
	v.Field("Str").Match("a+")
	v.Field("Bytes").Match("a+")
	result := v.Eval(obj)
	a.NotNil(result.Errors["Str"])
	a.NotNil(result.Errors["Bytes"])
	a.EqStr("not match with 'a+'", result.Errors["Str"][0].String())
	a.EqStr("not match with 'a+'", result.Errors["Bytes"][0].String())

	obj.Str = "bbb"
	obj.Bytes = []byte(obj.Str)
	result = v.Eval(obj)
	a.NotNil(result)
	a.NotNil(result.Errors["Str"])
	a.NotNil(result.Errors["Bytes"])
	a.EqStr("not match with 'a+'", result.Errors["Str"][0].String())
	a.EqStr("not match with 'a+'", result.Errors["Bytes"][0].String())

	obj.Str = "aaa"
	obj.Bytes = []byte(obj.Str)
	a.Nil(v.Eval(obj))
}

func TestUnmatch(t *testing.T) {
	a := assert.New(t)
	type TestObject struct {
		Str   string
		Bytes []byte
	}
	obj := &TestObject{}
	v := NewValidator()
	v.Field("Str").Unmatch("a+")
	v.Field("Bytes").Unmatch("a+")

	obj.Str = "aaa"
	obj.Bytes = []byte(obj.Str)
	result := v.Eval(obj)
	a.NotNil(result.Errors["Str"])
	a.NotNil(result.Errors["Bytes"])
	a.EqStr("match with 'a+'", result.Errors["Str"][0].String())
	a.EqStr("match with 'a+'", result.Errors["Bytes"][0].String())

	obj.Str = "bbb"
	obj.Bytes = []byte(obj.Str)
	a.Nil(v.Eval(obj))
}

func TestFunc(t *testing.T) {
	a := assert.New(t)
	type TestObject struct {
		Str string
	}
	obj := &TestObject{}
	v := NewValidator()
	v.Field("Str").Func(func(v interface{}) *FieldError {
		if v == "foo" {
			return NewFieldError(
				"Foo!!", nil,
			)
		}
		return nil
	})
	obj.Str = "foo"
	result := v.Eval(obj)
	a.NotNil(result.Errors["Str"])
	a.EqStr("Foo!!", result.Errors["Str"][0].String())

	obj.Str = "bbb"
	a.Nil(v.Eval(obj))
}
