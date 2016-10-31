package gaetest

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/speedland/go/x/xtesting/assert"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func mkTempfile(content string) string {
	f, err := ioutil.TempFile("", "gae-fixture-test")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(content)
	return f.Name()
}

type FixtureKind struct {
	IntValue      int
	FloatValue    float32
	BoolValue     bool
	StringValue   string
	DateTimeValue time.Time
	DateValue     time.Time
	BytesValue    []byte
	Slice         []string
	Struct        FixtureKindStruct
}

type FixtureKindStruct struct {
	Foo string
}

func TestDatastoreFixture(t *testing.T) {
	filepath := mkTempfile(`[{
    "_kind": "FixtureKind",
    "_key": "key1",
    "IntValue": 10,
    "FloatValue": 2.4,
    "BoolValue": true,
    "StringValue": "foobar",
    "BytesValue": "[]bytesfoobar",
    "DateTimeValue": "2014-01-02T14:02:50Z",
    "DateValue": "2014-01-02",
    "Slice": ["a", "b", "c"],
    "Struct": {
      "Foo": "bar"
    }
  },{
    "_kind": "FixtureKind",
    "_key": "key1",
    "_ns": "ns1",
    "StringValue": "withns1"
  }
]`)

	a := assert.New(t)
	var fk FixtureKind
	ctx := NewContext()
	a.Nil(DatastoreFixture(ctx, filepath, nil), "DatastoreFixture")

	key := datastore.NewKey(ctx, "FixtureKind", "key1", 0, nil)

	a.Nil(datastore.Get(ctx, key, &fk), "datastore.Get('key1') ")
	a.EqInt(10, fk.IntValue, "IntValue should be 10")
	a.EqFloat32(2.4, fk.FloatValue, "FloatValue should be 2.4")
	a.EqStr("foobar", fk.StringValue, "StringValue should be 'foobar'")
	a.EqStr("bytesfoobar", string(fk.BytesValue), "BytesValue should be 'foobar'")
	a.EqInt(3, len(fk.Slice), "len(Slice) should be 3")
	a.EqStr("a", string(fk.Slice[0]), "Slice[0] should be 'a'")
	a.EqStr("b", string(fk.Slice[1]), "Slice[0] should be 'a'")
	a.EqStr("c", string(fk.Slice[2]), "Slice[0] should be 'a'")

	a.EqTime(time.Date(2014, 01, 02, 14, 02, 50, 0, time.UTC), fk.DateTimeValue, "DateTimeValue should be 2014-01-02T14:02:50Z")
	a.EqTime(time.Date(2014, 01, 02, 0, 0, 0, 0, time.UTC), fk.DateValue, "DateTimeValue should be 2014-01-02T00:00:00Z")
	a.EqStr("bar", string(fk.Struct.Foo), "Struct.Foo should be 'bar'")

	// namespace
	ns1, _ := appengine.Namespace(ctx, "ns1")
	key = datastore.NewKey(ns1, "FixtureKind", "key1", 0, nil)
	a.Nil(datastore.Get(ns1, key, &fk), "datastore.Get('key1') /w ns1")
	a.EqStr("withns1", fk.StringValue, "StringValue should be 'withns1'")
}

func TestDatastoreFixtureWithBindings(t *testing.T) {
	filepath := mkTempfile(`[{
    "_kind": "FixtureKind",
    "_key": "key1",
    "IntValue": 10,
    "FloatValue": 2.4,
    "BoolValue": true,
    "StringValue": "foobar",
    "BytesValue": "[]bytesfoobar",
    "DateTimeValue": "{{now}}",
    "DateValue": "{{today}}"
  }]`)
	ctx := NewContext()
	a := assert.New(t)
	var fk FixtureKind
	DatastoreFixture(ctx, filepath, nil)
	key := datastore.NewKey(ctx, "FixtureKind", "key1", 0, nil)
	a.Nil(datastore.Get(ctx, key, &fk), "datastore.Get('key1') ")
}
