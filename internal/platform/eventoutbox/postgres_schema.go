package eventoutbox

import "fmt"

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
