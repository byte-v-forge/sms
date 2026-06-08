package eventoutbox

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/timex"
)

const (
	DefaultBatch          = 20
	DefaultInterval       = time.Second
	DefaultActiveInterval = 100 * time.Millisecond
)

type PendingProcessor interface {
	PublishPending(ctx context.Context, batch int) (int, error)
}

type WorkerConfig struct {
	Name           string
	Processor      PendingProcessor
	Batch          int
	Interval       time.Duration
	ActiveInterval time.Duration
	Logf           func(string, ...any)
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
		delay := cfg.Interval
		if published > 0 {
			delay = cfg.ActiveInterval
		}
		if err := timex.Sleep(ctx, delay); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}
	}
	return nil
}

func normalizeWorkerConfig(cfg WorkerConfig) WorkerConfig {
	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		cfg.Name = "event outbox"
	}
	if cfg.Batch <= 0 {
		cfg.Batch = DefaultBatch
	}
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultInterval
	}
	if cfg.ActiveInterval <= 0 {
		cfg.ActiveInterval = DefaultActiveInterval
	}
	if cfg.Logf == nil {
		cfg.Logf = log.Printf
	}
	return cfg
}
