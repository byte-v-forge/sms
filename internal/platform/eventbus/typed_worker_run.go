package eventbus

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

func RunTypedConsumerWorker[T proto.Message](ctx context.Context, cfg TypedConsumerWorkerConfig[T]) error {
	if cfg.Consumer == nil || cfg.NewMessage == nil || cfg.Handler == nil {
		return nil
	}
	cfg = normalizeTypedConsumerWorkerConfig(cfg)
	if err := cfg.Expected.ValidateEvent(cfg.NewMessage()); err != nil {
		return fmt.Errorf("configure %s expected event: %w", cfg.Name, err)
	}
	return RunConsumerWorker(ctx, ConsumerWorkerConfig{
		Name:            cfg.Name,
		Consumer:        cfg.Consumer,
		Batch:           cfg.Batch,
		FetchErrorDelay: cfg.FetchErrorDelay,
		Logf:            cfg.Logf,
		Handler: func(ctx context.Context, received ReceivedMessage) {
			handleTypedMessage(ctx, cfg, received)
		},
	})
}

func normalizeTypedConsumerWorkerConfig[T proto.Message](cfg TypedConsumerWorkerConfig[T]) TypedConsumerWorkerConfig[T] {
	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		cfg.Name = "typed event consumer"
	}
	cfg.MalformedLabel = strings.TrimSpace(cfg.MalformedLabel)
	if cfg.MalformedLabel == "" {
		cfg.MalformedLabel = fmt.Sprintf("terminate malformed %s", cfg.Name)
	}
	cfg.Logf = logger(cfg.Logf)
	return cfg
}
