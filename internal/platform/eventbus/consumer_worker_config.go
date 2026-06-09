package eventbus

import (
	"context"
	"strings"
	"time"
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
