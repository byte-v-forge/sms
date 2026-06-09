package eventoutbox

import (
	"context"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func publishMessage(ctx context.Context, publisher eventbus.Publisher, message eventbus.Message, options PublishOptions) error {
	publishCtx, cancel := context.WithTimeout(ctx, publishTimeout(options))
	_, err := publisher.Publish(publishCtx, message)
	cancel()
	return err
}
