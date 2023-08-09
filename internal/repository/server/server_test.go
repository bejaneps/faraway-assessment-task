package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/cache"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/db"
	"github.com/stretchr/testify/require"
)

var (
	testingDB    db.DB
	testingCache cache.Cache
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		return
	}

	var err error
	testingDB, err = initDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	testingCache, err = initCache()
	if err != nil {
		testingDB.Close()
		log.Fatal(err.Error())
	}

	exitCode := m.Run()
	testingDB.Close()
	testingCache.Close()
	os.Exit(exitCode) // doesn't respect defers
}

func initDB() (db.DB, error) {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		return nil, errors.New("DB_DSN is required for integration tests")
	}

	db, err := db.New(db.Config{
		DSN:  dbDSN,
		DBMS: db.PostgresDBMS,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return db, err
}

func initCache() (cache.Cache, error) {
	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		return nil, errors.New("CACHE_URL is required for integration tests")
	}

	cache, err := cache.New(cache.Config{
		URL: cacheURL,
		DB:  0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init cache: %w", err)
	}

	return cache, nil
}

func TestQuote(t *testing.T) {
	ctx := context.TODO()

	t.Run("success get quote from cache", func(t *testing.T) {
		err := testingCache.Seed(ctx)
		require.NoError(t, err)

		r := New(testingCache, testingDB)
		quote, err := r.Quote(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, quote)
	})

	t.Run("success get quote from database", func(t *testing.T) {
		err := testingDB.Seed(ctx)
		require.NoError(t, err)

		closedCache, err := initCache()
		require.NoError(t, err)
		require.NotNil(t, closedCache)
		require.NoError(t, closedCache.Close())

		r := New(closedCache, testingDB)
		quote, err := r.Quote(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, quote)
	})

	t.Run("fail both db and cache fail", func(t *testing.T) {
		closedCache, err := initCache()
		require.NoError(t, err)
		require.NotNil(t, closedCache)
		require.NoError(t, closedCache.Close())

		closedDB, err := initDB()
		require.NoError(t, err)
		require.NotNil(t, closedDB)
		require.NoError(t, closedDB.Close())

		r := New(closedCache, closedDB)
		quote, err := r.Quote(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get random quote from cache:")
		require.Contains(t, err.Error(), "failed to get random quote from database:")
		require.Empty(t, quote)
	})
}
