package db

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		return
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	testingDBPostgres, err = getPostgresClient(db)
	if err != nil {
		db.Close()
		log.Fatal(err.Error())
	}
	cleanupFuncPostgres()

	exitCode := m.Run()
	db.Close()
	os.Exit(exitCode)
}

func initDB() (DB, error) {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		return nil, errors.New("DB_DSN is required for integration tests")
	}

	db, err := New(Config{
		DSN:  dbDSN,
		DBMS: PostgresDBMS,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return db, err
}

func getPostgresClient(db DB) (*postgres, error) {
	var (
		p  = &postgres{}
		ok bool
	)

	p.conn, ok = db.client().(*dbr.Connection)
	if !ok {
		return nil, errors.New("wrong db client returned")
	}
	if p.conn == nil {
		return nil, errors.New("nil db client returned")
	}

	return p, nil
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d, err := initDB()
		require.NoError(t, err)
		require.NotNil(t, d)
		require.NoError(t, d.Close())
	})

	t.Run("fail unsupported dbms", func(t *testing.T) {
		d, err := New(Config{
			DBMS: "invalidDbms",
		})
		require.Error(t, err)
		require.Equal(t, err.Error(), "unsupported dbms invalidDbms, supported values are: postgres")
		require.Nil(t, d)
	})

	t.Run("fail db ping error", func(t *testing.T) {
		d, err := New(Config{
			DSN:  "invalidDsn",
			DBMS: PostgresDBMS,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to ping db")
		require.Nil(t, d)
	})
}
