package server

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/cache"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/db"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/models"
)

const (
	retryAttemptsCache = 3
	retryDelayCache    = 1 * time.Second
	retryAttemptsDB    = 2
	retryDelayDB       = 1 * time.Second
)

const selectRandomQuery = `SELECT quote FROM quotes ORDER BY random() LIMIT 1;`

// Quoter is used to interfact with quotes repository
type Quoter interface {
	// Quote gets random quote from cache or db
	Quote(ctx context.Context) (quote string, err error)
}

type Repo struct {
	c cache.Cache
	d db.DB
}

func New(c cache.Cache, d db.DB) *Repo {
	return &Repo{
		c: c,
		d: d,
	}
}

// Quote returns random quote from cache or db
func (r *Repo) Quote(ctx context.Context) (quote string, err error) {
	// first try getting from cache
	errCache := retry.Do(func() error {
		quote, err = r.c.SRandMember(ctx, models.CacheKeyQuotesSet)
		if err != nil {
			return err
		}

		return nil
	}, retry.Attempts(retryAttemptsCache), retry.Delay(retryDelayCache), retry.DelayType(retry.FixedDelay))
	if errCache == nil && quote != "" {
		return quote, nil
	}
	errCache = fmt.Errorf("failed to get random quote from cache: %w", errCache)

	// second try getting from database
	errDB := retry.Do(func() error {
		return r.d.Select(ctx, db.SelectArgs{
			Query:  selectRandomQuery,
			Result: &quote,
		})
	}, retry.Attempts(retryAttemptsDB), retry.Delay(retryDelayDB), retry.DelayType(retry.FixedDelay))
	if errDB == nil && quote != "" {
		return quote, nil
	}
	errDB = fmt.Errorf("failed to get random quote from database: %w", errDB)

	return "", fmt.Errorf("%w: %w", errCache, errDB)
}
