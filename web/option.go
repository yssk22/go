package web

import (
	"context"
	"crypto/sha256"
	"net/http"

	"github.com/speedland/go/x/xcrypto/xhmac"
)

// Option provies the option fields for web package.
type Option struct {
	// Option for hmac signature key, must not be nil. The default key is "speedland"
	HMACKey *xhmac.Base64
	// Option to initialize the request context. The default is nil.
	InitContext func(r *http.Request) context.Context
}

var DefaultOption = &Option{
	HMACKey: xhmac.NewBase64([]byte("speedland"), sha256.New),
}
