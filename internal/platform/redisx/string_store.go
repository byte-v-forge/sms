package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StringStore struct {
	client   redis.Cmdable
	keyspace Keyspace
	ttl      time.Duration
}

func NewStringStore(client redis.Cmdable, prefix string, ttl time.Duration) *StringStore {
	return &StringStore{
		client:   client,
		keyspace: NewKeyspace(prefix),
		ttl:      ttl,
	}
}

func (s *StringStore) DefaultTTL() time.Duration {
	if s == nil {
		return 0
	}
	return s.ttl
}

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

func (s *StringStore) SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return fmt.Errorf("redis string store key is required")
	}
	return s.client.Set(ctx, redisKey, value, s.effectiveTTL(ttl)).Err()
}

func (s *StringStore) redisKey(key string) (string, bool) {
	if s == nil || s.client == nil {
		return "", false
	}
	return s.keyspace.Key(key)
}

func (s *StringStore) effectiveTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return s.ttl
	}
	return ttl
}
