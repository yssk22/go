package keyvalue

import (
	"fmt"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

type MyMap map[string]interface{}

func (m MyMap) Get(key string) (interface{}, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return nil, KeyError(key)
}

func TestGetStringOr(t *testing.T) {
	a := assert.New(t)
	m := MyMap{
		"Foo": "1",
	}
	a.EqStr("1", GetStringOr(m, "Foo", "2"))
	a.EqStr("2", GetStringOr(m, "Bar", "2"))
}

func TestGetIntOr(t *testing.T) {
	a := assert.New(t)
	m := MyMap{
		"Foo":              1,
		"Int8":             int8(1),
		"StringNum":        "1",
		"InvalidStringNum": "a",
	}
	a.EqInt(1, GetIntOr(m, "Foo", 1))
	a.EqInt(1, GetIntOr(m, "Int8", 2))
	a.EqInt(1, GetIntOr(m, "StringNum", 2))
	a.EqInt(2, GetIntOr(m, "InvalidStringNum", 2))
}

func ExampleGetIntOr() {
	// type MyMap map[string]interface{}
	//
	// func (m MyMap) Get(key string) (interface{}, error) {
	// 	if v, ok := m[key]; ok {
	// 		return v, nil
	// 	}
	// 	return nil, KeyError(key)
	// }
	m := MyMap{
		"Foo":              1,
		"Int8":             int8(1),
		"StringNum":        "1",
		"InvalidStringNum": "a",
	}
	fmt.Println(GetIntOr(m, "Foo", 1))
	fmt.Println(GetIntOr(m, "Int8", 2))
	fmt.Println(GetIntOr(m, "StringNum", 2))
	fmt.Println(GetIntOr(m, "InvalidStringNum", 2))
	// Output:
	// 1
	// 1
	// 1
	// 2
}
