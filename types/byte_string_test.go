package types

import (
	"encoding/json"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestByteString_JSON(t *testing.T) {
	a := assert.New(t)
	type s struct {
		Value ByteString `json:"value"`
	}
	obj := &s{
		Value: []byte("value"),
	}
	buff, err := json.Marshal(obj)
	a.Nil(err)
	a.EqStr("{\"value\":\"value\"}", string(buff))

	var obj2 s
	a.Nil(json.Unmarshal(buff, &obj2))
	a.EqStr(string(obj.Value), string(obj2.Value))
}
