package cache

import (
	"context"
	"fmt"
)

// Cache interface for the client to use before accessing the datastore
type Cache interface {
	SetMulti(ctx context.Context, keys []string, values interface{}) error
	GetMulti(ctx context.Context, keys []string, dst interface{}) error
	DeleteMulti(ctx context.Context, keys []string) error
}

// ErrCacheKeyNotFound is an error alias
type ErrCacheKeyNotFound string

func (e ErrCacheKeyNotFound) Error() string {
	return fmt.Sprintf("key %q not found", string(e))
}
