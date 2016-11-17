package rgb

import (
	"encoding/json"
	"testing"
	"x/assert"
)

func TestRGB_Parse(t *testing.T) {
	a := assert.New(t)
	rgb, err := ParseRGB("#ff0000")
	a.Nil(err)
	a.EqInt(255<<16, int(rgb))

	rgb, err = ParseRGB("ff0000")
	a.Nil(err)
	a.EqInt(255<<16, int(rgb))

	rgb, err = ParseRGB("0xff0000")
	a.Nil(err)
	a.EqInt(255<<16, int(rgb))
}

func TestRGB_JSON(t *testing.T) {
	a := assert.New(t)
	type s struct {
		Value RGB `json:"value"`
	}
	obj := &s{RGB(255 << 16)}
	buff, err := json.Marshal(obj)
	a.Nil(err)
	a.EqStr("{\"value\":\"#ff0000\"}", string(buff))

	var obj2 s
	a.Nil(json.Unmarshal(buff, &obj2))
	a.EqInt(int(obj.Value), int(obj2.Value))
}
