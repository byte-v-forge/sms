package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (s *StringStore) Load(ctx context.Context, key string) (string, bool, error) {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return "", false, nil
	}
	value, err := s.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}
