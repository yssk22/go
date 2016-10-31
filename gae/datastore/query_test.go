package datastore

import (
	"fmt"
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestQuery_Filter(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	q := NewQuery("Example").Eq("ID", lazy.New("example-1"))
	_, err := q.GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(1, len(result))
}

func TestQuery_Order(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	q := NewQuery("Example").Desc("ID")
	_, err := q.GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(5, len(result))
	for i := range result {
		a.EqStr(fmt.Sprintf("example-%d", 5-i), result[i].ID)
	}
}
