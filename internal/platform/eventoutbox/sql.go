package eventoutbox

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidTableName = errors.New("event outbox table name is invalid")
	ErrNilDB            = errors.New("event outbox database handle is nil")
)

type PgxTx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func postgresIdentifier(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidTableName
	}
	parts := strings.Split(value, ".")
	for _, part := range parts {
		if !validPostgresIdentifierPart(part) {
			return "", fmt.Errorf("%w: %s", ErrInvalidTableName, value)
		}
	}
	return value, nil
}

func validPostgresIdentifierPart(value string) bool {
	if value == "" {
		return false
	}
	for i, r := range value {
		valid := r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || i > 0 && r >= '0' && r <= '9'
		if !valid {
			return false
		}
	}
	return true
}
