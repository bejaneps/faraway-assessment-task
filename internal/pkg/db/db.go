package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/lib/pq"
)

const (
	pqConnectionErrorClass   = "08"
	pqInterventionErrorClass = "57"

	startupErrDBStartingUp = "the database system is starting up"
	startupErrConnRefused  = "connection refused"
	startupErrFailedPing   = "failed to ping database"

	retryAttemptsDBPing = 5
	retryDelayDBPing    = 2 * time.Second
)

// DB is interface for dealing with database related queries
type DB interface {
	// Select selects record(s) from database and returns records with error if any
	Select(ctx context.Context, args SelectArgs) error
	// Close closes all db connections
	Close() error
	// Seed runs all seeds for database
	Seed(ctx context.Context) error

	// return underyling db client, used for tests
	client() interface{}
}

// SelectArgs is used for selecting records from database
type SelectArgs struct {
	// Query database query to run
	Query string
	// List of args to pass to above query
	Args []interface{}
	// Result field is used to decode result from database,
	// value for this field should be a pointer to underlying variable
	Result interface{}
}

type DBMS = string

const (
	PostgresDBMS = "postgres"
)

type Config struct {
	DSN                string        `env:"DB_DSN,required"`
	MaxOpenConnections int           `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"20"`
	MaxIdleConnections int           `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"10"`
	ConnMaxLifetime    time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
	DBMS               DBMS          `env:"DB_DBMS" envDefault:"postgres"`
}

func New(config Config) (DB, error) {
	if config.DBMS != PostgresDBMS {
		return nil, fmt.Errorf("unsupported dbms %s, supported values are: postgres", config.DBMS)
	}

	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// try connecting to db
	err = retry.Do(func() error {
		return db.Ping()
	}, retry.RetryIf(func(err error) bool {
		return isDBConnErr(err)
	}), retry.Attempts(retryAttemptsDBPing), retry.Delay(retryDelayDBPing), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	conn := &dbr.Connection{
		DB:            db,
		Dialect:       dialect.PostgreSQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}

	return &postgres{
		conn: conn,
	}, nil
}

func isDBConnErr(err error) bool {
	isPqErr := func(errp *pq.Error) bool {
		if string(errp.Code.Class()) == pqConnectionErrorClass ||
			string(errp.Code.Class()) == pqInterventionErrorClass {
			return true
		}

		return false
	}

	isExtStartupErr := func(err error) bool {
		if strings.Contains(err.Error(), startupErrDBStartingUp) ||
			strings.Contains(err.Error(), startupErrConnRefused) ||
			strings.Contains(err.Error(), startupErrFailedPing) {
			return true
		}

		return false
	}

	errp, ok := err.(*pq.Error)
	if ok && isPqErr(errp) || isExtStartupErr(err) {
		return true
	}

	return false
}
