package datastore

import (
	"fmt"
	"testing"

	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestQuery_Filter(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	q := NewQuery("Example").Eq("ID", "example-1")
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

func TestQuery_Limit(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	q := NewQuery("Example").Desc("ID").Limit(1)
	_, err := q.GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(1, len(result))
	a.EqStr("example-5", result[0].ID)
}

func TestQuery_KeysOnly(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	q := NewQuery("Example").Desc("ID").Limit(1).KeysOnly(true)
	keys, err := q.GetAll(gaetest.NewContext(), nil)
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("/Example,example-5", keys[0].String())

	q = q.KeysOnly(false)
	keys, err = q.GetAll(gaetest.NewContext(), nil)
	a.EqStr("datastore: invalid entity type", err.Error())
}
