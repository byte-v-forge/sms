package eventoutbox

import (
	"context"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func publishRow(ctx context.Context, publisher eventbus.Publisher, row Row, updates Updates, options PublishOptions) (bool, error) {
	message, err := MessageFromEnvelope(row.Envelope)
	if err != nil {
		return discardRow(ctx, row, updates, options, err)
	}
	if err := publishMessage(ctx, publisher, message, options); err != nil {
		return retryRow(ctx, row, updates, options, err)
	}
	publishedAt := optionNow(options).Unix()
	return true, updates.MarkPublished(ctx, row.EventID, publishedAt)
}
