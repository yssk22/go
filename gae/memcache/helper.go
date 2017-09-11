package memcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/speedland/go/lazy"

	"github.com/speedland/go/x/xlog"

	"github.com/speedland/go/web"

	"github.com/speedland/go/web/response"

	"context"
)

// CachedObjectWithExpiration execute generator function and store the result with `key` and set the value into dst with expiration.
func CachedObjectWithExpiration(
	ctx context.Context,
	key string,
	expiration time.Duration,
	dst interface{},
	generator func() (interface{}, error),
	force bool,
) error {
	if !force {
		if err := Get(ctx, key, dst); err == nil {
			return nil
		}
	}
	gen, err := generator()
	if err != nil {
		return err
	}
	if err = SetWithExpire(ctx, key, gen, expiration); err != nil {
		return err
	}
	buff, _ := json.Marshal(gen)
	json.Unmarshal(buff, dst)
	return nil
}

var cachableHeaderKeys = map[string]bool{
	response.ContentType: true,
}

func RegisterCachableHeaderKeys(keys ...string) {
	for _, k := range keys {
		cachableHeaderKeys[strings.ToLower(k)] = true
	}
}

type cachedResponseBody []byte

func (c cachedResponseBody) Render(ctx context.Context, w io.Writer) {
	w.Write([]byte(c))
}

const maxCachableBodySize = 1 * 1024 * 900 // 900KB as Memcache supports 1MB

// CacheResponse wraps web.Handler to support cache
func CacheResponse(name lazy.Value, h web.Handler) web.Handler {
	return CacheResponseWithExpire(name, 0, h)
}

// CacheResponseWithExpire wraps web.Handler to support cache
func CacheResponseWithExpire(name lazy.Value, expire time.Duration, h web.Handler) web.Handler {
	// var ckey = fmt.Sprintf("response-%s", name)
	type cachedResponse struct {
		Status response.HTTPStatus `json:"s"`
		Header http.Header         `json:"h"`
		Body   cachedResponseBody  `json:"b"`
	}
	var ckey = fmt.Sprintf("cache-response-%s", name)
	var prefix = ckey + " "
	return web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		logger := xlog.WithContext(ctx).WithPrefix(prefix)
		var cr cachedResponse
		if err := Get(ctx, ckey, &cr); err == nil {
			if cr.Status != 0 {
				return &response.Response{
					Status: cr.Status,
					Header: cr.Header,
					Body:   cr.Body,
				}
			}
			logger.Warnf("discard cache as it's invalid: (status: %d, size: %d)", cr.Status, len(cr.Body))
		}

		resp := h.Process(req, next)
		if len(resp.Cookies) > 0 {
			logger.Warnf("disabling cache capability: including cookies")
			return resp
		}
		for k, v := range resp.Header {
			if _, ok := cachableHeaderKeys[strings.ToLower(k)]; !ok {
				logger.Warnf("disabling cache capability: including unregistered headres (header-key: %s, header-value: %v)", ckey, k, v)
				return resp
			}
		}
		var buff bytes.Buffer
		resp.Body.Render(ctx, &buff)
		cr = cachedResponse{
			Status: resp.Status,
			Header: resp.Header,
			Body:   cachedResponseBody(buff.Bytes()),
		}
		resp.Body = cr.Body
		if sz := len(cr.Body); sz > maxCachableBodySize {
			logger.Warnf("disabling cache capability: too big to cache (cache-key:%s, body size: %v)", ckey, sz)
			return resp
		}
		if err := SetWithExpire(ctx, ckey, &cr, expire); err != nil {
			logger.Errorf("could not use response cache (key: %s): %v", ckey, err)
		}
		return resp
	})
}
