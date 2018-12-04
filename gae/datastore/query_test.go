package datastore

import (
	"fmt"
	"testing"

	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestQuery_KeysOnly(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	keys, err := NewQuery("Example").Eq("ID", "example-1").KeysOnly().GetAll(gaetest.NewContext(), nil)
	a.Nil(err)
	a.EqInt(1, len(keys))
}

func TestQuery_Filter(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	_, err := NewQuery("Example").Eq("ID", "example-1").GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(1, len(result))
}

func TestQuery_Order(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	_, err := NewQuery("Example").Desc("ID").GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(5, len(result))
	for i := range result {
		a.EqStr(fmt.Sprintf("example-%d", 5-i), result[i].ID)
	}
}

func TestQuery_Limit(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	_, err := NewQuery("Example").Desc("ID").Limit(1).GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(1, len(result))
	a.EqStr("example-5", result[0].ID)
}
