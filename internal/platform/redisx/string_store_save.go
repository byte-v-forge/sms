package redisx

import (
	"context"
	"fmt"
	"time"
)

func (s *StringStore) SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return fmt.Errorf("redis string store key is required")
	}
	return s.client.Set(ctx, redisKey, value, s.effectiveTTL(ttl)).Err()
}
