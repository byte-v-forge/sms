package app

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newRequiredPostgresPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("SMS_PG_DSN or PG_DSN is required")
	}
	return pgxpool.New(ctx, dsn)
}
