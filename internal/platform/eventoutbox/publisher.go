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
