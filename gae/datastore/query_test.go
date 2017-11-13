package datastore

import (
	"fmt"
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xtesting/assert"
)

const queryLoggerKey = "github.com/speedland/go/gae/datastore.Query"

func TestQuery_Filter(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	q := NewQuery("Example", queryLoggerKey).Eq("ID", lazy.New("example-1"))
	_, err := q.GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(1, len(result))
}

func TestQuery_Order(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	var result []Example
	q := NewQuery("Example", queryLoggerKey).Desc("ID")
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
	q := NewQuery("Example", queryLoggerKey).Desc("ID").Limit(lazy.New(1))
	_, err := q.GetAll(gaetest.NewContext(), &result)
	a.Nil(err)
	a.EqInt(1, len(result))
	a.EqStr("example-5", result[0].ID)
}

func TestQuery_KeysOnly(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestQuery.json", nil))

	q := NewQuery("Example", queryLoggerKey).Desc("ID").Limit(lazy.New(1)).KeysOnly(true)
	keys, err := q.GetAll(gaetest.NewContext(), nil)
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("/Example,example-5", keys[0].String())

	q = q.KeysOnly(false)
	keys, err = q.GetAll(gaetest.NewContext(), nil)
	a.EqStr("datastore: invalid entity type", err.Error())
}
