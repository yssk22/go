package keyvalue

import (
	"fmt"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestGetStringOr(t *testing.T) {
	a := assert.New(t)
	m := Map{
		"Foo":  "1",
		"Bar":  []string{"a"},
		"Hoge": []string{},
	}
	a.EqStr("1", GetStringOr(m, "Foo", "2"))
	a.EqStr("a", GetStringOr(m, "Bar", "2"))
	a.EqStr("2", GetStringOr(m, "Hoge", "2"))
	a.EqStr("2", GetStringOr(m, "Me", "2"))
}

func TestGetIntOr(t *testing.T) {
	a := assert.New(t)
	m := Map{
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
	m := Map{
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
