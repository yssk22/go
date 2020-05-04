package datastore

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/yssk22/go/x/xtesting/assert"
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

func TestTestEnv(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	key := datastore.NameKey("MyKind", "foo", nil)
	client := testEnv.NewClient()
	defer client.Close()
	a.Nil(testEnv.memcache.SetMulti(ctx, []string{"foo"}, []string{"bar"}))
	_, err := client.inner.Put(ctx, key, &Example{ID: "bar"})
	a.Nil(err)

	var e Example
	s := make([]string, 1, 1)
	a.Nil(testEnv.memcache.GetMulti(ctx, []string{"foo"}, s))
	a.Nil(client.inner.Get(ctx, key, &e))

	a.EqStr(s[0], "bar")
	a.EqStr(e.ID, "bar")
	a.Nil(testEnv.Reset())

	a.NotNil(testEnv.memcache.GetMulti(ctx, []string{"foo"}, s))
	c, err := client.inner.Count(ctx, datastore.NewQuery("MyKind"))
	a.Nil(err)
	a.EqInt(0, c)
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
	a.Nil(testEnv.LoadFixture(filepath))
	client := testEnv.NewClient()
	defer client.Close()

	var fk FixtureKind
	key := datastore.NameKey("FixtureKind", "key1", nil)
	a.Nil(client.inner.Get(context.Background(), key, &fk), "client.Get('key1') ")

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
	key = datastore.NameKey("FixtureKind", "key1", nil)
	key.Namespace = "ns1"
	a.Nil(client.inner.Get(context.Background(), key, &fk), "client.Get('ns1.key1') ")
	a.EqStr("withns1", fk.StringValue, "StringValue should be 'withns1'")
}
