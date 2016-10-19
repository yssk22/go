package iterator

import (
	"fmt"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func ExampleParallelSlice() {
	var a = []int{0, 1, 2, 3, 4}
	ParallelSlice(a, DefaultParallelOption, func(i int, v int) error {
		a[i] = a[i] + 1
		return nil
	})
	fmt.Println(a)
	// Output: [1 2 3 4 5]
}

func ExampleParallelSlice_Struct() {
	type T struct {
		i int
	}
	var a = []T{
		T{0}, T{1}, T{2}, T{3},
	}
	ParallelSlice(a, DefaultParallelOption, func(i int, t *T) error {
		t.i = t.i + 1
		return nil
	})
	fmt.Println(a)
	// Output: [{0} {1} {2} {3}]
}

func TestParallelSlice_withMaxMaxConcurrency(t *testing.T) {
	assert := assert.New(t)
	var opts = &ParallelOption{
		MaxConcurrency: 3,
	}
	var a = []int{0, 1, 2, 3, 4}
	assert.Nil(
		ParallelSlice(a, opts, func(i int, v int) error {
			assert.EqInt(i, v)
			a[i] = v + 1
			return nil
		}),
	)
	for i := range a {
		assert.EqInt(i+1, a[i])
	}
}
