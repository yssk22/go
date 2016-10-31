package slice

import (
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestToInterfaceSlice(t *testing.T) {
	a := assert.New(t)
	type T struct {
		i int
	}

	// built in
	intSlice := []int{1, 2, 3}
	ia := ToInterfaceSlice(intSlice)
	a.EqInt(2, ia[1].(int))

	// []T
	tSlice := []T{T{1}, T{2}, T{3}}
	ia = ToInterfaceSlice(tSlice)
	a.EqInt(2, ia[1].(T).i)

	// []*T
	ptrTSlice := []*T{&T{1}, &T{2}, &T{3}}
	ia = ToInterfaceSlice(ptrTSlice)
	a.EqInt(2, ia[1].(*T).i)

}
