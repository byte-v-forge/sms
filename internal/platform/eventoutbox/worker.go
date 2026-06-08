package eventoutbox

import (
	"context"
	"errors"

	"github.com/byte-v-forge/sms/internal/platform/timex"
)

type PendingProcessor interface {
	PublishPending(ctx context.Context, batch int) (int, error)
}

func RunWorker(ctx context.Context, cfg WorkerConfig) error {
	if cfg.Processor == nil {
		return nil
	}
	cfg = normalizeWorkerConfig(cfg)
	for ctx.Err() == nil {
		published, err := cfg.Processor.PublishPending(ctx, cfg.Batch)
		if err != nil {
			cfg.Logf("publish %s failed: %v", cfg.Name, err)
		}
		if err := timex.Sleep(ctx, workerDelay(cfg, published)); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}
	}
	return nil
}
