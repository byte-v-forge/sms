package hotstreamnats

import (
	"context"
	"errors"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func (b *Bus) Subscribe(ctx context.Context, filter hotstream.Filter) (*hotstream.Subscription, error) {
	if b == nil || b.hub == nil {
		return nil, errors.New("hotstream bus is not configured")
	}
	return b.hub.Subscribe(ctx, filter)
}
