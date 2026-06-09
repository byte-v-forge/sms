package eventoutbox

import (
	"context"
	"fmt"
)

func (u pgxUpdates) MarkDiscarded(ctx context.Context, eventID string, lastError string, updatedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = $1, last_error = $2, updated_at = $3
		WHERE event_id = $4
	`, u.table), StatusDiscarded, lastError, updatedAt, eventID)
	return err
}
