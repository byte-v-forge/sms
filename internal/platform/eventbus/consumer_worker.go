package eventbus

import (
	"context"
)

func RunConsumerWorker(ctx context.Context, cfg ConsumerWorkerConfig) error {
	if cfg.Consumer == nil || cfg.Handler == nil {
		return nil
	}
	cfg = normalizeConsumerWorkerConfig(cfg)
	for ctx.Err() == nil {
		messages, err := cfg.Consumer.Fetch(ctx, cfg.Batch)
		if err != nil {
			if done, err := handleConsumerFetchError(ctx, cfg, err); done || err != nil {
				return err
			}
			continue
		}
		for _, message := range messages {
			cfg.Handler(ctx, message)
		}
	}
	return nil
}
