package slice

import (
	"fmt"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func ExampleSplitByLength() {
	var a = []int{0, 1, 2, 3, 4}
	var b = SplitByLength(a, 2).([][]int)
	fmt.Println(b)
	// Output: [[0 1] [2 3] [4]]
}

func Test_SplitByLength(t *testing.T) {
	a := assert.New(t)
	var a1 = []int{1, 2, 3, 4, 5}
	var a2 = SplitByLength(a1, 3).([][]int)
	a.EqInt(1, a2[0][0])
	a.EqInt(2, a2[0][1])
	a.EqInt(3, a2[0][2])
	a.EqInt(2, len(a2[1]))
	a.EqInt(4, a2[1][0])
	a.EqInt(5, a2[1][1])
}

func Test_ToInterface(t *testing.T) {
	a := assert.New(t)
	type T struct {
		i int
	}

	// built in
	intSlice := []int{1, 2, 3}
	ia := ToInterface(intSlice)
	a.EqInt(2, ia[1].(int))

	// []T
	tSlice := []T{T{1}, T{2}, T{3}}
	ia = ToInterface(tSlice)
	a.EqInt(2, ia[1].(T).i)

	// []*T
	ptrTSlice := []*T{&T{1}, &T{2}, &T{3}}
	ia = ToInterface(ptrTSlice)
	a.EqInt(2, ia[1].(*T).i)
}

func Test_ToAddr(t *testing.T) {
	a := assert.New(t)
	type Foo struct {
		field string
	}
	list := ToAddr([]Foo{Foo{
		field: "foo",
	}}).([]*Foo)
	a.EqStr("foo", list[0].field)
}

func Test_ToElem(t *testing.T) {
	a := assert.New(t)
	type Foo struct {
		field string
	}
	list := ToElem([]*Foo{&Foo{
		field: "foo",
	}}).([]Foo)
	a.EqStr("foo", list[0].field)
}
