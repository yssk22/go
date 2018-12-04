package counter

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yssk22/go/x/xerrors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

// Config is a counter shard configuration
// @datastore kind=CounterConfig
type Config struct {
	Key       string    `json:"key" ent:"id"`
	NumShards int       `json:"num_shareds"`
	UpdatedAt time.Time `json:"updated_at" ent:"timestamp"`
}

// Shard is a counter shard.
// @datastore kind=CounterShard
type Shard struct {
	Key        string    `json:"key" ent:"id"`
	CounterKey string    `json:"counter_key"`
	Count      int       `json:"count"`
	UpdatedAt  time.Time `json:"updated_at" ent:"timestamp"`
}

const defaultNumShareds = 10
const countMemcacheExpiration = 60 * time.Second

func countMemcacheKey(key string) string {
	return fmt.Sprintf("counter.%s", key)
}

// Count returns the count value of the given key.
func Count(ctx context.Context, key string) (int, error) {
	var total int
	memkey := countMemcacheKey(key)
	if _, err := memcache.JSON.Get(ctx, memkey, &total); err == nil {
		return total, nil
	}
	_, shards, err := NewShardQuery().EqCounterKey(key).GetAll(ctx)
	if err != nil {
		return 0, err
	}
	for _, s := range shards {
		total += s.Count
	}
	memcache.JSON.Set(ctx, &memcache.Item{
		Key:        memkey,
		Object:     &total,
		Expiration: countMemcacheExpiration,
	})
	return total, nil
}

// MustCount is like Count but panics if an error occurrs.
func MustCount(ctx context.Context, key string) int {
	c, err := Count(ctx, key)
	xerrors.MustNil(err)
	return c
}

// Reset resets the counter of the given key
func Reset(ctx context.Context, key string) error {
	kind := NewShardKind()
	keys, _, err := NewShardQuery().EqCounterKey(key).GetAll(ctx)
	if err != nil {
		return err
	}
	if _, err := kind.DeleteMulti(ctx, keys); err != nil {
		return err
	}
	memcache.Delete(ctx, countMemcacheKey(key))
	return nil
}

// MustReset is like Reset but panics if an error occurrs.
func MustReset(ctx context.Context, key string) {
	xerrors.MustNil(Reset(ctx, key))
}

// Increment increments the counter of the given key
func Increment(ctx context.Context, key string) error {
	configKind := NewConfigKind()
	shardKind := NewShardKind()
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var cfg *Config
		_, cfg, err := configKind.Get(ctx, key)
		if err != nil {
			return err
		}
		if cfg == nil {
			cfg = &Config{
				Key:       key,
				NumShards: defaultNumShareds,
			}
			_, err := configKind.Put(ctx, cfg)
			if err != nil {
				return err
			}
		}
		shardKey := fmt.Sprintf("%s.shard.%d", key, rand.Intn(cfg.NumShards))
		_, shard, err := shardKind.Get(ctx, shardKey)
		if err != nil {
			return err
		}
		if shard == nil {
			shard = &Shard{
				Key:        shardKey,
				CounterKey: key,
				Count:      0,
			}
		}
		shard.Count++
		_, err = shardKind.Put(ctx, shard)
		return err
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		return xerrors.Wrap(err, "could not increment %s", key)
	}
	memcache.IncrementExisting(ctx, countMemcacheKey(key), 1)
	return nil
}

// MustIncrement is like Increment but panics if an error occurrs.
func MustIncrement(ctx context.Context, key string) {
	xerrors.MustNil(Increment(ctx, key))
}
