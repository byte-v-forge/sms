package eventoutbox

import (
	"context"
	"fmt"
)

func (u pgxUpdates) MarkRetry(ctx context.Context, eventID string, attemptCount int32, nextAttemptAt int64, lastError string, updatedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET attempt_count = $1, next_attempt_at = $2, last_error = $3, updated_at = $4
		WHERE event_id = $5
	`, u.table), attemptCount, nextAttemptAt, lastError, updatedAt, eventID)
	return err
}
