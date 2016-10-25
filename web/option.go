package web

import (
	"crypto/sha256"

	"github.com/speedland/go/x/xcrypto/xhmac"
)

// Option provies the option fields for web package.
type Option struct {
	// Option for hmac signature key, must not be nil. The default key is "speedland"
	HMACKey *xhmac.Base64
}

var DefaultOption = &Option{
	HMACKey: xhmac.NewBase64([]byte("speedland"), sha256.New),
}
