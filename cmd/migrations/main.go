package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/cache"
	_ "github.com/bejaneps/faraway-assessment-task/internal/pkg/debug" // prints debug info
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/db"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/golang-migrate/migrate/v4"
)

const (
	migrationsPath = "file://./"

	retryAttempts = 10
	retryDelay    = 2 * time.Second
)

func main() {
	log.Info("starting up migrations")

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("DB_DSN is required for migrations")
	}

	dbClient, err := db.New(db.Config{
		DSN:  dbDSN,
		DBMS: db.PostgresDBMS,
	})
	if err != nil {
		log.Fatal("failed to init db", log.StdError(err))
	}
	defer dbClient.Close()

	m, err := migrate.New(migrationsPath, dbDSN)
	if err != nil {
		log.Fatal("failed to init db with database instance", log.StdError(err))
	}

	if err := m.Up(); err != nil {
		log.Fatal("failed to migrate up", log.StdError(err))
	}

	// run seeds
	if err = dbClient.Seed(context.Background()); err == nil {
		runCacheSeeds()
		log.Info("migrations and seeds ran successfully")
		return
	}

	log.Error("failed to run db seeds", log.StdError(err))

	// run db seeds
	err = retry.Do(
		func() error {
			return m.Steps(-1)
		},
		retry.Attempts(retryAttempts),
		retry.Delay(retryDelay),
		retry.DelayType(retry.FixedDelay),
	)
	if err != nil {
		log.Fatal("failed to migrate down", log.StdError(err))
	}

	log.Info("migrations and seed reverted successfully")
}

func runCacheSeeds() {
	// run cache seeds, no need to migrate down if they fail, it's just cache
	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		log.Fatal("CACHE_URL is required for migrations")
	}

	cacheDB := os.Getenv("CACHE_DB")
	if cacheDB == "" {
		log.Fatal("CACHE_DB is required for migrations")
	}

	cacheDBInt, err := strconv.Atoi(cacheDB)
	if err != nil {
		log.Fatal("wrong number for CACHE_DB", log.String("number", cacheDB))
	}

	cacheClient, err := cache.New(cache.Config{
		URL: cacheURL,
		DB:  cacheDBInt,
	})
	if err != nil {
		log.Error("failed to init cache", log.StdError(err), log.String("url", cacheURL))
		return
	}
	defer cacheClient.Close()

	if err := cacheClient.Seed(context.Background()); err != nil {
		log.Error("failed to run cache seeds", log.StdError(err))
		return
	}
}
