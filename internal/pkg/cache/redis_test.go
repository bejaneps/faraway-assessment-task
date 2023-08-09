package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/models"
	pkgQuotes "github.com/bejaneps/faraway-assessment-task/internal/pkg/quotes"
	"github.com/stretchr/testify/require"
)

var (
	cleanupFuncRedis = func() {
		testingCacheRedis.c.FlushAll(context.TODO())
	}

	testingCacheRedis *redis
)

func TestSRandMember_Redis(t *testing.T) {
	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanupFuncRedis)

		err := testingCacheRedis.Seed(ctx)
		require.NoError(t, err)

		quote, err := testingCacheRedis.SRandMember(ctx, models.CacheKeyQuotesSet)
		require.NoError(t, err)
		require.NotEmpty(t, quote)

		ok := false
		for _, seed := range pkgQuotes.Seeds {
			if quote == seed.Quote {
				ok = true
				break
			}
		}

		require.True(t, ok, "failed to find random quote from cache")
	})

	t.Run("fail cache error", func(t *testing.T) {
		c, err := initCache()
		require.NoError(t, err)
		r, err := getRedisClient(c)
		require.NoError(t, err)
		r.Close()

		quote, err := r.SRandMember(ctx, models.CacheKeyQuotesSet)
		require.Error(t, err)
		require.Empty(t, quote)
	})
}

func TestSeed_Redis(t *testing.T) {
	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanupFuncRedis)

		err := testingCacheRedis.Seed(ctx)
		require.NoError(t, err)

		// check all keys were stored in db
		quotes, err := testingCacheRedis.c.Keys(ctx, "quotes#*").Result()
		require.NoError(t, err)
		require.Equal(t, len(pkgQuotes.Seeds), len(quotes))

		for i := 0; i < len(pkgQuotes.Seeds); i++ {
			key := fmt.Sprintf(models.CacheKeyQuotes, pkgQuotes.Seeds[i].ID)
			quote, err := testingCacheRedis.c.Get(ctx, key).Result()
			require.NoError(t, err)
			require.Equal(t, pkgQuotes.Seeds[i].Quote, quote)
		}

		// check keys were appended to set
		quotesSet, err := testingCacheRedis.c.SMembers(ctx, models.CacheKeyQuotesSet).Result()
		require.NoError(t, err)
		require.Equal(t, len(pkgQuotes.Seeds), len(quotesSet))
	})

	t.Run("success run seeds twice", func(t *testing.T) {
		t.Cleanup(cleanupFuncRedis)

		err := testingCacheRedis.Seed(ctx)
		require.NoError(t, err)

		// shouldn't add new keys to cache
		err = testingCacheRedis.Seed(ctx)
		require.NoError(t, err)

		quotes, err := testingCacheRedis.c.Keys(ctx, "quotes#*").Result()
		require.NoError(t, err)
		require.Equal(t, len(pkgQuotes.Seeds), len(quotes))

		for i := 0; i < len(pkgQuotes.Seeds); i++ {
			key := fmt.Sprintf(models.CacheKeyQuotes, pkgQuotes.Seeds[i].ID)
			quote, err := testingCacheRedis.c.Get(ctx, key).Result()
			require.NoError(t, err)
			require.Equal(t, pkgQuotes.Seeds[i].Quote, quote)
		}

		// check keys were appended to set only once
		quotesSet, err := testingCacheRedis.c.SMembers(ctx, models.CacheKeyQuotesSet).Result()
		require.NoError(t, err)
		require.Equal(t, len(pkgQuotes.Seeds), len(quotesSet))
	})

	t.Run("fail cache error", func(t *testing.T) {
		c, err := initCache()
		require.NoError(t, err)
		r, err := getRedisClient(c)
		require.NoError(t, err)
		r.Close()

		err = r.Seed(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to run seed")
	})
}
