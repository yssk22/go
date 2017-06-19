package xos

import (
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func Test_TryStats(t *testing.T) {
	a := assert.New(t)
	info, err := TryStats("./not_exists.txt", "./fixtures/stats_a.txt", "./fixtures/stats_b.txt")
	a.Nil(err)
	a.EqStr("stats_a.txt", info.Name())

	info, err = TryStats("./not_exists.txt", "./fixtures/stats_b.txt", "./fixtures/stats_a.txt")
	a.Nil(err)
	a.EqStr("stats_b.txt", info.Name())

	info, err = TryStats("./not_exists.txt", "./not_exist_b.txt")
	a.NotNil(err)
}
