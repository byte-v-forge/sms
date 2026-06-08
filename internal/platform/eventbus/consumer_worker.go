package eventbus

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/timex"
)

const DefaultFetchErrorDelay = time.Second

type LogFunc func(string, ...any)

type MessageHandler func(context.Context, ReceivedMessage)

type ConsumerWorkerConfig struct {
	Name            string
	Consumer        Consumer
	Handler         MessageHandler
	Batch           int
	FetchErrorDelay time.Duration
	Logf            LogFunc
}

func RunConsumerWorker(ctx context.Context, cfg ConsumerWorkerConfig) error {
	if cfg.Consumer == nil || cfg.Handler == nil {
		return nil
	}
	cfg = normalizeConsumerWorkerConfig(cfg)
	for ctx.Err() == nil {
		messages, err := cfg.Consumer.Fetch(ctx, cfg.Batch)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			cfg.Logf("fetch %s failed: %v", cfg.Name, err)
			if err := timex.Sleep(ctx, cfg.FetchErrorDelay); err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
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

func normalizeConsumerWorkerConfig(cfg ConsumerWorkerConfig) ConsumerWorkerConfig {
	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		cfg.Name = "event consumer"
	}
	if cfg.FetchErrorDelay <= 0 {
		cfg.FetchErrorDelay = DefaultFetchErrorDelay
	}
	cfg.Logf = logger(cfg.Logf)
	return cfg
}
