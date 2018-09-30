package xhmac

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestBase64_SignAndUnsign(t *testing.T) {
	a := assert.New(t)
	message := "mizukifukumura"
	key := "morningmusume"
	hmac := NewBase64([]byte(key), nil)
	signed := hmac.Sign([]byte(message))
	unsigned, err := hmac.Unsign(signed)
	a.Nil(err)
	a.EqStr(message, string(unsigned))
}

func TestBase64_Unsign_ModifiedError(t *testing.T) {
	a := assert.New(t)
	message := "mizukifukumura"
	key := "morningmusume"
	hmac := NewBase64([]byte(key), nil)
	signed := hmac.Sign([]byte(message))
	signed[0] = byte('a')
	_, err := hmac.Unsign(signed)
	a.NotNil(err)
	a.OK(ErrSignatureModified == err)
}
