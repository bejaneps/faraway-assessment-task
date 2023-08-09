package db

import (
	"context"
	"testing"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/models"
	pkgQuotes "github.com/bejaneps/faraway-assessment-task/internal/pkg/quotes"
	"github.com/stretchr/testify/require"
)

var (
	cleanupFuncPostgres = func() {
		session := testingDBPostgres.conn.NewSession(nil)
		session.ExecContext(context.TODO(), "DELETE FROM quotes WHERE id IN (1, 2, 3, 4);")
	}

	testingDBPostgres *postgres
)

func TestSelect_Postgres(t *testing.T) {
	ctx := context.TODO()

	seed := prepareSeedsPostgresFunc()[0]

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanupFuncPostgres)

		session := testingDBPostgres.conn.NewSession(nil)
		_, err := session.ExecContext(ctx, seed)
		require.NoError(t, err)

		var quote models.TableQuotes
		err = testingDBPostgres.Select(ctx, SelectArgs{
			Query:  "SELECT * FROM quotes WHERE id = ?;",
			Args:   []interface{}{1},
			Result: &quote,
		})
		require.NoError(t, err)
		require.Equal(t, pkgQuotes.Seeds[0].ID, quote.ID)
		require.Equal(t, pkgQuotes.Seeds[0].Quote, quote.Quote)
	})

	t.Run("fail invalid sql or error from db", func(t *testing.T) {
		t.Cleanup(cleanupFuncPostgres)

		session := testingDBPostgres.conn.NewSession(nil)
		_, err := session.ExecContext(ctx, seed)
		require.NoError(t, err)

		var quote models.TableQuotes
		err = testingDBPostgres.Select(ctx, SelectArgs{
			Query:  "invalidSql",
			Args:   []interface{}{1},
			Result: &quote,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to select record(s) from db:")
	})
}

func TestSeed_Postgres(t *testing.T) {
	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanupFuncPostgres)

		err := testingDBPostgres.Seed(ctx)
		require.NoError(t, err)

		session := testingDBPostgres.conn.NewSession(nil)

		var quotes []models.TableQuotes
		_, err = session.Select("*").From(models.TableNameQuotes).OrderAsc("id").LoadContext(ctx, &quotes)
		require.NoError(t, err)
		require.Equal(t, len(pkgQuotes.Seeds), len(quotes))
		require.Equal(t, pkgQuotes.Seeds[0].ID, quotes[0].ID)
		require.Equal(t, pkgQuotes.Seeds[0].Quote, quotes[0].Quote)
	})

	t.Run("success run seeds twice", func(t *testing.T) {
		t.Cleanup(cleanupFuncPostgres)

		err := testingDBPostgres.Seed(ctx)
		require.NoError(t, err)

		// shouldn't add new rows to db
		err = testingDBPostgres.Seed(ctx)
		require.NoError(t, err)

		session := testingDBPostgres.conn.NewSession(nil)

		var quotes []models.TableQuotes
		_, err = session.Select("*").From(models.TableNameQuotes).OrderAsc("id").LoadContext(ctx, &quotes)
		require.NoError(t, err)
		require.Equal(t, len(pkgQuotes.Seeds), len(quotes))
		require.Equal(t, pkgQuotes.Seeds[0].ID, quotes[0].ID)
		require.Equal(t, pkgQuotes.Seeds[0].Quote, quotes[0].Quote)
	})

	t.Run("fail begin tx error", func(t *testing.T) {
		db, err := initDB()
		require.NoError(t, err)
		p, err := getPostgresClient(db)
		require.NoError(t, err)
		p.Close()

		err = p.Seed(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to begin tx for running seeds:")
	})

	//TODO: this test sometimes fail,
	// because it interferes with other tests as they are running tests on same db instance
	t.Run("fail invalid seed error", func(t *testing.T) {
		cleanupFuncPostgres()
		prepareSeedsPostgresFunc = func() []string {
			return []string{
				"INSERT INTO quotes(id, quote) VALUES (1, 'test1')",
				"invalidSql",
			}
		}
		t.Cleanup(func() {
			prepareSeedsPostgresFunc = prepareSeedsPostgres
		})

		err := testingDBPostgres.Seed(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to run seed")

		// check if transaction was rolled back
		session := testingDBPostgres.conn.NewSession(nil)

		var quotes []models.TableQuotes
		_, err = session.Select("*").From(models.TableNameQuotes).OrderAsc("id").LoadContext(ctx, &quotes)
		require.NoError(t, err)
		require.Empty(t, quotes)
	})
}
