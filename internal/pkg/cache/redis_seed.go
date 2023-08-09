package cache

import (
	"context"
	"fmt"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/models"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/quotes"
)

func (r *redis) Seed(ctx context.Context) error {
	var seedsToAdd []string

	for _, seed := range quotes.Seeds {
		ok, err := r.c.SetNX(
			ctx, fmt.Sprintf(models.CacheKeyQuotes, seed.ID), seed.Quote, 0,
		).Result()
		if err != nil {
			return fmt.Errorf("failed to run seed (%d: %s): %w", seed.ID, seed.Quote, err)
		}

		// add seed to array, so we can get random number from array later
		if !ok {
			continue
		}

		seedsToAdd = append(seedsToAdd, seed.Quote)
	}

	if len(seedsToAdd) == 0 {
		return nil
	}

	if err := r.c.SAdd(ctx, models.CacheKeyQuotesSet, seedsToAdd).Err(); err != nil {
		return fmt.Errorf("failed to add seeds to set: %w", err)
	}

	return nil
}
