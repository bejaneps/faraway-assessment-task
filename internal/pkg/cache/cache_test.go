package cache

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	pkgRedis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		return
	}

	cache, err := initCache()
	if err != nil {
		log.Fatal(err.Error())
	}

	testingCacheRedis, err = getRedisClient(cache)
	if err != nil {
		cache.Close()
		log.Fatal(err.Error())
	}
	cleanupFuncRedis()

	exitCode := m.Run()
	cache.Close()
	os.Exit(exitCode)
}

func initCache() (Cache, error) {
	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		return nil, errors.New("CACHE_URL is required for integration tests")
	}

	cache, err := New(Config{
		URL: cacheURL,
		DB:  0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init cache: %w", err)
	}

	return cache, nil
}

func getRedisClient(c Cache) (*redis, error) {
	var (
		r  = &redis{}
		ok bool
	)

	r.c, ok = c.client().(*pkgRedis.Client)
	if !ok {
		return nil, errors.New("wrong cache client returned")
	}
	if r.c == nil {
		return nil, errors.New("nil cache client returned")
	}

	return r, nil
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, err := initCache()
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NoError(t, c.Close())
	})

	t.Run("fail cache ping error", func(t *testing.T) {
		prevEnv := os.Getenv("CACHE_URL")
		t.Cleanup(func() {
			os.Setenv("CACHE_URL", prevEnv)
		})
		require.NoError(t, os.Setenv("CACHE_URL", "invalid:6379"))
		c, err := initCache()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to ping cache")
		require.Nil(t, c)
	})
}
