package eventoutbox

import (
	"context"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func PublishRows(ctx context.Context, publisher eventbus.Publisher, rows []Row, updates Updates, options PublishOptions) (int, error) {
	if publisher == nil {
		return 0, ErrNilPublisher
	}
	if updates == nil {
		return 0, ErrNilUpdates
	}
	published := 0
	for _, row := range rows {
		if ctx.Err() != nil {
			return published, ctx.Err()
		}
		rowPublished, err := publishRow(ctx, publisher, row, updates, options)
		if err != nil {
			return published, err
		}
		if rowPublished {
			published++
		}
	}
	return published, nil
}

func publishRow(ctx context.Context, publisher eventbus.Publisher, row Row, updates Updates, options PublishOptions) (bool, error) {
	message, err := MessageFromEnvelope(row.Envelope)
	if err != nil {
		discardedAt := optionNow(options).Unix()
		return false, updates.MarkDiscarded(ctx, row.EventID, TruncateError(err), discardedAt)
	}

	publishCtx, cancel := context.WithTimeout(ctx, publishTimeout(options))
	_, err = publisher.Publish(publishCtx, message)
	cancel()
	if err != nil {
		nextAttempt := row.AttemptCount + 1
		updatedAt := optionNow(options)
		nextAttemptAt := updatedAt.Add(retryDelay(options, nextAttempt)).Unix()
		return false, updates.MarkRetry(ctx, row.EventID, nextAttempt, nextAttemptAt, TruncateError(err), updatedAt.Unix())
	}
	publishedAt := optionNow(options).Unix()
	return true, updates.MarkPublished(ctx, row.EventID, publishedAt)
}
