package db

import (
	"context"
	"fmt"

	"github.com/gocraft/dbr/v2"
)

type postgres struct {
	conn *dbr.Connection
}

// Select selects record(s) from Postgres
func (p *postgres) Select(ctx context.Context, args SelectArgs) error {
	res := args.Result

	session := p.conn.NewSession(nil)
	_, err := session.SelectBySql(args.Query, args.Args...).LoadContext(ctx, res)
	if err != nil {
		return fmt.Errorf("failed to select record(s) from db: %w", err)
	}

	return nil
}

func (p *postgres) Close() error {
	return p.conn.Close()
}

func (p *postgres) client() interface{} {
	return p.conn
}
