package xhttp

import (
	"net/http"
	"testing"

	"github.com/speedland/go/x/xcrypto/xhmac"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestSignCookie(t *testing.T) {
	a := assert.New(t)
	hmac := xhmac.NewBase64([]byte("mykey"), nil)
	c := &http.Cookie{
		Name:  "foo",
		Value: "bar",
	}
	cc := SignCookie(c, hmac)
	a.EqStr(cc.Name, c.Name)
	a.EqStr(cc.Value, hmac.SignString(c.Value))
}

func TestUnsignCookie(t *testing.T) {
	a := assert.New(t)
	hmac := xhmac.NewBase64([]byte("mykey"), nil)
	cc := &http.Cookie{
		Name:  "foo",
		Value: hmac.SignString("value"),
	}
	c, err := UnsignCookie(cc, hmac)
	a.Nil(err)
	a.EqStr(c.Name, cc.Name)
	a.EqStr(c.Value, "value")
}
