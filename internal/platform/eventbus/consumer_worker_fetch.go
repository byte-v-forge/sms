package eventbus

import (
	"context"
	"errors"

	"github.com/byte-v-forge/sms/internal/platform/timex"
)

func handleConsumerFetchError(ctx context.Context, cfg ConsumerWorkerConfig, err error) (bool, error) {
	if ctx.Err() != nil {
		return true, nil
	}
	cfg.Logf("fetch %s failed: %v", cfg.Name, err)
	if err := timex.Sleep(ctx, cfg.FetchErrorDelay); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return true, nil
		}
		return true, err
	}
	return false, nil
}
