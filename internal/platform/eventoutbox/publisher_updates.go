package eventoutbox

import "context"

func discardRow(ctx context.Context, row Row, updates Updates, options PublishOptions, err error) (bool, error) {
	discardedAt := optionNow(options).Unix()
	return false, updates.MarkDiscarded(ctx, row.EventID, TruncateError(err), discardedAt)
}

func retryRow(ctx context.Context, row Row, updates Updates, options PublishOptions, err error) (bool, error) {
	nextAttempt := row.AttemptCount + 1
	updatedAt := optionNow(options)
	nextAttemptAt := updatedAt.Add(retryDelay(options, nextAttempt)).Unix()
	return false, updates.MarkRetry(ctx, row.EventID, nextAttempt, nextAttemptAt, TruncateError(err), updatedAt.Unix())
}
