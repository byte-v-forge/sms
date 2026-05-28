package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func validatePostgresTables(ctx context.Context, pool *pgxpool.Pool, tables ...string) error {
	for _, table := range tables {
		var exists bool
		if err := pool.QueryRow(ctx, "SELECT to_regclass($1) IS NOT NULL", "public."+table).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("database schema is not migrated: missing table %s", table)
		}
	}
	return nil
}
