package slice

import (
	"fmt"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
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
	fmt.Println(a2)
	a.EqInt(1, a2[0][0])
	a.EqInt(2, a2[0][1])
	a.EqInt(3, a2[0][2])
	a.EqInt(2, len(a2[1]))
	a.EqInt(4, a2[1][0])
	a.EqInt(5, a2[1][1])
}
