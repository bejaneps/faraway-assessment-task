package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	redisPkg "github.com/redis/go-redis/v9"
)

const (
	retryAttemptsCachePing = 5
	retryDelayCachePing    = 2 * time.Second
)

// Cache is interface for dealing with cache related queries
type Cache interface {
	// SRandMember executes SRANDMEMBER command in cache
	SRandMember(ctx context.Context, key string) (string, error)
	// Close closes all cache connections
	Close() error
	// Seed runs all seeds for cache
	Seed(ctx context.Context) error

	// return underyling cache client, used for tests
	client() interface{}
}

type Config struct {
	URL string `env:"CACHE_URL,required"`
	DB  int    `env:"CACHE_DB,required"`
}

func New(config Config) (Cache, error) {
	r := &redis{
		c: redisPkg.NewClient(&redisPkg.Options{
			Addr: config.URL,
			DB:   config.DB,
		}),
	}

	err := retry.Do(func() error {
		_, err := r.c.Ping(context.TODO()).Result()
		return err
	}, retry.Attempts(retryAttemptsCachePing), retry.Delay(retryDelayCachePing), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return nil, fmt.Errorf("failed to ping cache: %w", err)
	}

	return r, nil
}
