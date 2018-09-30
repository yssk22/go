package web

import (
	"context"
	"crypto/sha256"
	"net/http"

	"github.com/yssk22/go/x/xcrypto/xhmac"
)

// Option provies the option fields for web package.
type Option struct {
	// Option for hmac signature key, must not be nil. The default key is "github.com/yssk22"
	HMACKey *xhmac.Base64
	// Option to initialize the request context. The default is nil.
	InitContext func(r *http.Request) context.Context
}

var DefaultOption = &Option{
	HMACKey: xhmac.NewBase64([]byte("github.com/yssk22"), sha256.New),
}
