package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/yssk22/go/validator"
	"github.com/yssk22/go/x/xtime"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_ParameterParser_Get_Primitives(t *testing.T) {
	type TestStruct struct {
		BoolVal  bool      `json:"bool_val"`
		IntVal   int       `json:"int_val"`
		FloatVal float64   `json:"float_val"`
		StrVal   string    `json:"str_val"`
		TimeVal  time.Time `json:"time_val"`
	}
	var params = TestStruct{}

	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("bool_val", RequestParameterFieldTypeBool)
	pp.Type("int_val", RequestParameterFieldTypeInt)
	pp.Type("float_val", RequestParameterFieldTypeFloat)
	pp.Type("str_val", RequestParameterFieldTypeString)
	pp.Type("time_val", RequestParameterFieldTypeTime)

	req, _ := http.NewRequest("GET", "/?bool_val=true&int_val=1&float_val=1.2&str_val=a&time_val=2019/01/01", nil)
	res := pp.Parse(req, &params)
	a.Nil(res)
	a.OK(params.BoolVal)
	a.EqInt(1, params.IntVal)
	a.EqFloat64(1.2, params.FloatVal)
	a.EqStr("a", params.StrVal)
	a.EqTime(xtime.MustParseDateDefault("2019/01/01"), params.TimeVal)

	params = TestStruct{}
	req, _ = http.NewRequest("GET", "/", nil)
	res = pp.Parse(req, &params)
	a.Nil(res)
	a.OK(!params.BoolVal)
	a.EqInt(0, params.IntVal)
	a.EqFloat64(0.0, params.FloatVal)
	a.EqStr("", params.StrVal)
	a.OK(params.TimeVal.IsZero())
}

func Test_ParameterParser_Get_PrimitivesWithDefault(t *testing.T) {
	type TestStruct struct {
		BoolVal  bool      `json:"bool_val"`
		IntVal   int       `json:"int_val"`
		FloatVal float64   `json:"float_val"`
		StrVal   string    `json:"str_val"`
		TimeVal  time.Time `json:"time_val"`
	}
	var params = TestStruct{}

	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("bool_val", RequestParameterFieldTypeBool).Default("bool_val", true)
	pp.Type("int_val", RequestParameterFieldTypeInt).Default("int_val", 10)
	pp.Type("float_val", RequestParameterFieldTypeFloat).Default("float_val", 2.3)
	pp.Type("str_val", RequestParameterFieldTypeString).Default("str_val", "abc")
	pp.Type("time_val", RequestParameterFieldTypeTime).Default("time_val", xtime.MustParseDateDefault("1996/10/30"))

	req, _ := http.NewRequest("GET", "/?bool_val=true&int_val=1&float_val=1.2&str_val=a&time_val=2019/01/01", nil)
	res := pp.Parse(req, &params)
	a.Nil(res)
	a.OK(params.BoolVal)
	a.EqInt(1, params.IntVal)
	a.EqFloat64(1.2, params.FloatVal)
	a.EqStr("a", params.StrVal)
	a.EqTime(xtime.MustParseDateDefault("2019/01/01"), params.TimeVal)

	params = TestStruct{}
	req, _ = http.NewRequest("GET", "/", nil)
	res = pp.Parse(req, &params)
	a.Nil(res)
	a.OK(params.BoolVal)
	a.EqInt(10, params.IntVal)
	a.EqFloat64(2.3, params.FloatVal)
	a.EqStr("abc", params.StrVal)
	a.EqTime(xtime.MustParseDateDefault("1996/10/30"), params.TimeVal)
}

func Test_ParameterParser_Get_PtrPrimitives(t *testing.T) {
	type TestStruct struct {
		BoolVal  *bool      `json:"bool_val"`
		IntVal   *int       `json:"int_val"`
		FloatVal *float64   `json:"float_val"`
		StrVal   *string    `json:"str_val"`
		TimeVal  *time.Time `json:"time_val"`
	}
	var params = TestStruct{}

	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("bool_val", RequestParameterFieldTypeBool)
	pp.Type("int_val", RequestParameterFieldTypeInt)
	pp.Type("float_val", RequestParameterFieldTypeFloat)
	pp.Type("str_val", RequestParameterFieldTypeString)
	pp.Type("time_val", RequestParameterFieldTypeTime)

	req, _ := http.NewRequest("GET", "/?bool_val=true&int_val=1&float_val=1.2&str_val=a&time_val=2019/01/01", nil)
	res := pp.Parse(req, &params)
	a.Nil(res)
	a.OK(*params.BoolVal)
	a.EqInt(1, *params.IntVal)
	a.EqFloat64(1.2, *params.FloatVal)
	a.EqStr("a", *params.StrVal)
	a.EqTime(xtime.MustParseDateDefault("2019/01/01"), *params.TimeVal)

	params = TestStruct{}
	req, _ = http.NewRequest("GET", "/", nil)
	res = pp.Parse(req, &params)
	a.Nil(res)
	a.Nil(params.BoolVal)
	a.Nil(params.IntVal)
	a.Nil(params.FloatVal)
	a.Nil(params.StrVal)
	a.Nil(params.TimeVal)
}

func Test_ParameterParser_Get_PtrPrimitivesWithDefault(t *testing.T) {
	type TestStruct struct {
		BoolVal  *bool      `json:"bool_val"`
		IntVal   *int       `json:"int_val"`
		FloatVal *float64   `json:"float_val"`
		StrVal   *string    `json:"str_val"`
		TimeVal  *time.Time `json:"time_val"`
	}

	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("bool_val", RequestParameterFieldTypeBool).Default("bool_val", true)
	pp.Type("int_val", RequestParameterFieldTypeInt).Default("int_val", 10)
	pp.Type("float_val", RequestParameterFieldTypeFloat).Default("float_val", 2.3)
	pp.Type("str_val", RequestParameterFieldTypeString).Default("str_val", "abc")
	pp.Type("time_val", RequestParameterFieldTypeTime).Default("time_val", xtime.MustParseDateDefault("1996/10/30"))

	var params = TestStruct{}
	req, _ := http.NewRequest("GET", "/?bool_val=true&int_val=1&float_val=1.2&str_val=a&time_val=2019/01/01", nil)
	res := pp.Parse(req, &params)
	a.Nil(res)
	a.OK(*params.BoolVal)
	a.EqInt(1, *params.IntVal)
	a.EqFloat64(1.2, *params.FloatVal)
	a.EqStr("a", *params.StrVal)
	a.EqTime(xtime.MustParseDateDefault("2019/01/01"), *params.TimeVal)

	params = TestStruct{}
	req, _ = http.NewRequest("GET", "/", nil)
	res = pp.Parse(req, &params)
	a.Nil(res)
	a.OK(*params.BoolVal)
	a.EqInt(10, *params.IntVal)
	a.EqFloat64(2.3, *params.FloatVal)
	a.EqStr("abc", *params.StrVal)
	a.EqTime(xtime.MustParseDateDefault("1996/10/30"), *params.TimeVal)
}

func Test_ParameterParser_Get_Array(t *testing.T) {
	type TestStruct struct {
		IntVal []int `json:"int_val"`
	}
	var params = TestStruct{}

	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("int_val", RequestParameterFieldTypeArray)

	req, _ := http.NewRequest("GET", "/?int_val=[1,2,3]", nil)
	res := pp.Parse(req, &params)
	a.Nil(res)
	a.EqInt(1, params.IntVal[0])
	a.EqInt(2, params.IntVal[1])
	a.EqInt(3, params.IntVal[2])

	params = TestStruct{}
	req, _ = http.NewRequest("GET", "/", nil)
	res = pp.Parse(req, &params)
	a.Nil(res)
	a.Nil(params.IntVal)
}

func Test_ParameterParser_Get_Object(t *testing.T) {
	type TestInner2 struct {
		FloatVal float64   `json:"float_val"`
		TimeVal  time.Time `json:"time_val"`
	}
	type TestInner struct {
		IntVal  int         `json:"int_val"`
		TimeVal time.Time   `json:"time_val"`
		Inner2  *TestInner2 `json:"inner2"`
	}
	type TestStruct struct {
		IntVal int        `json:"int_val"`
		Inner  *TestInner `json:"inner"`
	}
	var expected = &TestStruct{
		Inner: &TestInner{
			IntVal:  10,
			TimeVal: xtime.MustParseDateDefault("2019/10/30"),
			Inner2: &TestInner2{
				FloatVal: 23.1,
				TimeVal:  xtime.MustParseDateDefault("2019/02/22"),
			},
		},
	}
	buff, _ := json.Marshal(expected.Inner)

	var params = TestStruct{}

	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("inner", RequestParameterFieldTypeObject)

	log.Println(fmt.Sprintf("/?int_val=1&inner=%s", string(buff)))
	req, _ := http.NewRequest("GET", fmt.Sprintf("/?inner=%s", url.QueryEscape(string(buff))), nil)
	res := pp.Parse(req, &params)
	a.Nil(res)
	a.EqInt(expected.IntVal, params.IntVal)
	a.EqInt(expected.Inner.IntVal, params.Inner.IntVal)
	a.EqTime(expected.Inner.TimeVal, params.Inner.TimeVal)
	a.EqFloat64(expected.Inner.Inner2.FloatVal, params.Inner.Inner2.FloatVal)
	a.EqTime(expected.Inner.Inner2.TimeVal, params.Inner.Inner2.TimeVal)

	params = TestStruct{}
	req, _ = http.NewRequest("GET", "/", nil)
	res = pp.Parse(req, &params)
	a.Nil(res)
	a.Nil(params.Inner)
}

func Test_ParameterParser_Get_Required(t *testing.T) {
	type TestStruct struct {
		BoolVal *bool `json:"bool_val"`
	}
	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("bool_val", RequestParameterFieldTypeBool).Required("bool_val")
	var params = TestStruct{}
	req, _ := http.NewRequest("GET", "/", nil)
	err := pp.Parse(req, &params)
	a.NotNil(err)
	a.EqStr("required", getFieldErrors(err)["bool_val"][0])

	req, _ = http.NewRequest("GET", "/?bool_val=true", nil)
	err = pp.Parse(req, &params)
	a.Nil(err)
}

type TestStructValidatable struct {
	IntVal int `json:"int_val"`
}

func (t *TestStructValidatable) Validate(ctx context.Context, errors FieldErrorCollection) error {
	errors.Add("int_val", validator.Int().Min(2).Validate(t.IntVal))
	return nil
}

func Test_ParameterParser_Get_Validatable(t *testing.T) {
	a := assert.New(t)
	pp := NewParameterParser(RequestParameterFormatQuery)
	pp.Type("int_val", RequestParameterFieldTypeInt).Required("int_val")
	var params = TestStructValidatable{}
	req, _ := http.NewRequest("GET", "/?int_val=1", nil)
	err := pp.Parse(req, &params)
	a.NotNil(err)
	a.EqStr("must be more than or equal to 2", getFieldErrors(err)["int_val"][0])

	req, _ = http.NewRequest("GET", "/?int_val=2", nil)
	err = pp.Parse(req, &params)
	a.Nil(err)
}

func getFieldErrors(err *Error) FieldErrorCollection {
	return err.Extra.(map[string]interface{})["errors"].(FieldErrorCollection)
}
