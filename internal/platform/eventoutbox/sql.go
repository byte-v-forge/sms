package eventoutbox

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrInvalidTableName = errors.New("event outbox table name is invalid")
	ErrNilDB            = errors.New("event outbox database handle is nil")
)

type PgxTx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func PostgresSchemaStatements(table string, pendingIndex string) ([]string, error) {
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	indexName, err := postgresIdentifier(pendingIndex)
	if err != nil {
		return nil, err
	}
	return []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			event_id TEXT PRIMARY KEY,
			subject TEXT NOT NULL,
			event_name TEXT NOT NULL,
			idempotency_key TEXT NOT NULL DEFAULT '',
			envelope BYTEA NOT NULL,
			status TEXT NOT NULL DEFAULT 'PENDING',
			attempt_count INT NOT NULL DEFAULT 0,
			next_attempt_at BIGINT NOT NULL DEFAULT 0,
			last_error TEXT NOT NULL DEFAULT '',
			published_at BIGINT NOT NULL DEFAULT 0,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		)`, tableName),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s (status, next_attempt_at, created_at)`, indexName, tableName),
	}, nil
}

func InsertRecordPgx(ctx context.Context, tx PgxTx, table string, record Record, now int64) error {
	if tx == nil {
		return ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return err
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	_, err = tx.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			event_id, subject, event_name, idempotency_key, envelope, status,
			attempt_count, next_attempt_at, last_error, published_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, 0, 0, '', 0, $7, $7)
		ON CONFLICT (event_id) DO NOTHING
	`, tableName), record.EventID, record.Subject, record.EventName, record.IdempotencyKey, record.Envelope, StatusPending, now)
	return err
}

func ClaimPendingPgx(ctx context.Context, tx PgxTx, table string, batch int, now int64) ([]Row, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	if batch <= 0 {
		batch = DefaultBatch
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	rows, err := tx.Query(ctx, fmt.Sprintf(`
		SELECT event_id, envelope, attempt_count
		FROM %s
		WHERE status = $1 AND next_attempt_at <= $2
		ORDER BY created_at ASC
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`, tableName), StatusPending, now, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Row{}
	for rows.Next() {
		var row Row
		if err := rows.Scan(&row.EventID, &row.Envelope, &row.AttemptCount); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func NewPgxUpdates(tx PgxTx, table string) (Updates, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	return pgxUpdates{tx: tx, table: tableName}, nil
}

type pgxUpdates struct {
	tx    PgxTx
	table string
}

func (u pgxUpdates) MarkPublished(ctx context.Context, eventID string, publishedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = $1, published_at = $2, updated_at = $2, last_error = ''
		WHERE event_id = $3
	`, u.table), StatusPublished, publishedAt, eventID)
	return err
}

func (u pgxUpdates) MarkRetry(ctx context.Context, eventID string, attemptCount int32, nextAttemptAt int64, lastError string, updatedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET attempt_count = $1, next_attempt_at = $2, last_error = $3, updated_at = $4
		WHERE event_id = $5
	`, u.table), attemptCount, nextAttemptAt, lastError, updatedAt, eventID)
	return err
}

func (u pgxUpdates) MarkDiscarded(ctx context.Context, eventID string, lastError string, updatedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = $1, last_error = $2, updated_at = $3
		WHERE event_id = $4
	`, u.table), StatusDiscarded, lastError, updatedAt, eventID)
	return err
}

func InsertRecordGORM(ctx context.Context, tx *gorm.DB, table string, record Record, now int64) error {
	if tx == nil {
		return ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return err
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	return tx.WithContext(ctx).Exec(fmt.Sprintf(`
		INSERT INTO %s (
			event_id, subject, event_name, idempotency_key, envelope, status,
			attempt_count, next_attempt_at, last_error, published_at, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, 0, 0, '', 0, ?, ?)
		ON CONFLICT (event_id) DO NOTHING
	`, tableName), record.EventID, record.Subject, record.EventName, record.IdempotencyKey, record.Envelope, StatusPending, now, now).Error
}

func ClaimPendingGORM(ctx context.Context, tx *gorm.DB, table string, batch int, now int64) ([]Row, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	if batch <= 0 {
		batch = DefaultBatch
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	rows := []Row{}
	err = tx.WithContext(ctx).Raw(fmt.Sprintf(`
		SELECT event_id, envelope, attempt_count
		FROM %s
		WHERE status = ? AND next_attempt_at <= ?
		ORDER BY created_at ASC
		LIMIT ?
		FOR UPDATE SKIP LOCKED
	`, tableName), StatusPending, now, batch).Scan(&rows).Error
	return rows, err
}

func NewGORMUpdates(tx *gorm.DB, table string) (Updates, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	return gormUpdates{tx: tx, table: tableName}, nil
}

type gormUpdates struct {
	tx    *gorm.DB
	table string
}

func (u gormUpdates) MarkPublished(ctx context.Context, eventID string, publishedAt int64) error {
	return u.tx.WithContext(ctx).Exec(fmt.Sprintf(`
		UPDATE %s
		SET status = ?, published_at = ?, updated_at = ?, last_error = ''
		WHERE event_id = ?
	`, u.table), StatusPublished, publishedAt, publishedAt, eventID).Error
}

func (u gormUpdates) MarkRetry(ctx context.Context, eventID string, attemptCount int32, nextAttemptAt int64, lastError string, updatedAt int64) error {
	return u.tx.WithContext(ctx).Exec(fmt.Sprintf(`
		UPDATE %s
		SET attempt_count = ?, next_attempt_at = ?, last_error = ?, updated_at = ?
		WHERE event_id = ?
	`, u.table), attemptCount, nextAttemptAt, lastError, updatedAt, eventID).Error
}

func (u gormUpdates) MarkDiscarded(ctx context.Context, eventID string, lastError string, updatedAt int64) error {
	return u.tx.WithContext(ctx).Exec(fmt.Sprintf(`
		UPDATE %s
		SET status = ?, last_error = ?, updated_at = ?
		WHERE event_id = ?
	`, u.table), StatusDiscarded, lastError, updatedAt, eventID).Error
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
