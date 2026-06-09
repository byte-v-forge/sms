package eventoutbox

import (
	"context"
	"fmt"
)

func (u pgxUpdates) MarkPublished(ctx context.Context, eventID string, publishedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = $1, published_at = $2, updated_at = $2, last_error = ''
		WHERE event_id = $3
	`, u.table), StatusPublished, publishedAt, eventID)
	return err
}
