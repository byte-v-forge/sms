package app

import (
	"context"
	"time"
)

type TTLStringStore interface {
	DefaultTTL() time.Duration
	Load(ctx context.Context, key string) (string, bool, error)
	SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error
}
