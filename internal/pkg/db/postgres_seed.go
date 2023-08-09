package db

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	pkgQuotes "github.com/bejaneps/faraway-assessment-task/internal/pkg/quotes"
)

const (
	retryAttemptsPostgresSeed = 2
	retryDelayPostgresSeed    = 1 * time.Second
)

const queryInsertSeedPostgres = `
INSERT INTO quotes(id, quote)
VALUES (%d, '%s')
ON CONFLICT DO NOTHING;
`

func (p *postgres) Seed(ctx context.Context) error {
	tx, err := p.conn.NewSession(nil).BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx for running seeds: %w", err)
	}

	retryOpts := []retry.Option{
		retry.Attempts(retryAttemptsPostgresSeed),
		retry.Delay(retryDelayPostgresSeed),
		retry.DelayType(retry.FixedDelay),
	}
	seeds := prepareSeedsPostgresFunc()
	for _, seed := range seeds { // execute each seed query
		err := retry.Do(func() error {
			_, err := tx.ExecContext(ctx, seed)
			return err
		}, retryOpts...)
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = fmt.Errorf("failed to rollback tx for running seeds: %w", errRollback)
			}

			return fmt.Errorf("failed to run seed (%s): %w", seed, err)
		}
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit tx for running seeds: %w", err)
		if errRollback := tx.Rollback(); errRollback != nil {
			err = fmt.Errorf("failed to rollback tx for running seeds: %w", errRollback)
		}
	}

	return err
}

// used for tests
var prepareSeedsPostgresFunc = prepareSeedsPostgres

func prepareSeedsPostgres() []string {
	seeds := make([]string, 0, len(pkgQuotes.Seeds))
	for _, seed := range pkgQuotes.Seeds {
		seeds = append(seeds, fmt.Sprintf(queryInsertSeedPostgres, seed.ID, seed.Quote))
	}

	return seeds
}
