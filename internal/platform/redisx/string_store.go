package redisx

import (
	"context"
	"fmt"
	"strings"
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

func (s *StringStore) LoadMany(ctx context.Context, keys ...string) (map[string]string, error) {
	cleanKeys, redisKeys := s.redisKeys(keys)
	if len(redisKeys) == 0 {
		return map[string]string{}, nil
	}
	values, err := s.client.MGet(ctx, redisKeys...).Result()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(values))
	for idx, raw := range values {
		value, ok := redisStringValue(raw)
		if !ok {
			continue
		}
		out[cleanKeys[idx]] = value
	}
	return out, nil
}

func (s *StringStore) Save(ctx context.Context, key string, value string) error {
	return s.SaveTTL(ctx, key, value, s.ttl)
}

func (s *StringStore) SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return fmt.Errorf("redis string store key is required")
	}
	ttl = s.effectiveTTL(ttl)
	return s.client.Set(ctx, redisKey, value, ttl).Err()
}

func (s *StringStore) Delete(ctx context.Context, key string) error {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return nil
	}
	return s.client.Del(ctx, redisKey).Err()
}

func (s *StringStore) HashLoadMany(ctx context.Context, key string, fields ...string) (map[string]string, error) {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return map[string]string{}, nil
	}
	cleanFields := cleanHashFields(fields)
	if len(cleanFields) == 0 {
		return map[string]string{}, nil
	}
	values, err := s.client.HMGet(ctx, redisKey, cleanFields...).Result()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(values))
	for idx, raw := range values {
		value, ok := redisStringValue(raw)
		if !ok {
			continue
		}
		out[cleanFields[idx]] = value
	}
	return out, nil
}

func (s *StringStore) HashSaveTTL(ctx context.Context, key string, values map[string]string, ttl time.Duration) error {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return fmt.Errorf("redis hash store key is required")
	}
	cleanValues := cleanHashValues(values)
	if len(cleanValues) == 0 {
		return nil
	}
	if err := s.client.HSet(ctx, redisKey, cleanValues).Err(); err != nil {
		return err
	}
	ttl = s.effectiveTTL(ttl)
	if ttl <= 0 {
		return nil
	}
	current, err := s.client.TTL(ctx, redisKey).Result()
	if err != nil {
		return err
	}
	if current <= 0 || ttl > current {
		return s.client.Expire(ctx, redisKey, ttl).Err()
	}
	return nil
}

func (s *StringStore) HashDelete(ctx context.Context, key string, fields ...string) error {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return nil
	}
	cleanFields := cleanHashFields(fields)
	if len(cleanFields) == 0 {
		return nil
	}
	return s.client.HDel(ctx, redisKey, cleanFields...).Err()
}

func (s *StringStore) redisKeys(keys []string) ([]string, []string) {
	seen := map[string]struct{}{}
	cleanKeys := make([]string, 0, len(keys))
	redisKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		cleanKey, redisKey, ok := s.cleanRedisKey(key)
		if !ok {
			continue
		}
		if _, exists := seen[cleanKey]; exists {
			continue
		}
		seen[cleanKey] = struct{}{}
		cleanKeys = append(cleanKeys, cleanKey)
		redisKeys = append(redisKeys, redisKey)
	}
	return cleanKeys, redisKeys
}

func (s *StringStore) redisKey(key string) (string, bool) {
	_, redisKey, ok := s.cleanRedisKey(key)
	return redisKey, ok
}

func (s *StringStore) cleanRedisKey(key string) (string, string, bool) {
	if s == nil || s.client == nil {
		return "", "", false
	}
	return s.keyspace.CleanKey(key)
}

func (s *StringStore) effectiveTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return s.ttl
	}
	return ttl
}

func redisStringValue(raw any) (string, bool) {
	switch value := raw.(type) {
	case nil:
		return "", false
	case string:
		return value, true
	case []byte:
		return string(value), true
	default:
		return fmt.Sprint(value), true
	}
}

func cleanHashFields(fields []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		if _, exists := seen[field]; exists {
			continue
		}
		seen[field] = struct{}{}
		out = append(out, field)
	}
	return out
}

func cleanHashValues(values map[string]string) map[string]string {
	out := make(map[string]string, len(values))
	for field, value := range values {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		out[field] = value
	}
	return out
}
