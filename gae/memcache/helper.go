package memcache

import (
	"encoding/json"
	"time"

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
