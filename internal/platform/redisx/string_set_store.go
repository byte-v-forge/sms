package redisx

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

type StringSetStore struct {
	client   redis.Cmdable
	keyspace Keyspace
}

func NewStringSetStore(client redis.Cmdable, prefix string) *StringSetStore {
	return &StringSetStore{client: client, keyspace: NewKeyspace(prefix)}
}

func (s *StringSetStore) Add(ctx context.Context, key string, members ...string) error {
	redisKey, cleanMembers, ok := s.redisKeyAndMembers(key, members)
	if !ok {
		return nil
	}
	return s.client.SAdd(ctx, redisKey, stringArgs(cleanMembers)...).Err()
}

func (s *StringSetStore) Remove(ctx context.Context, key string, members ...string) error {
	redisKey, cleanMembers, ok := s.redisKeyAndMembers(key, members)
	if !ok {
		return nil
	}
	return s.client.SRem(ctx, redisKey, stringArgs(cleanMembers)...).Err()
}

func (s *StringSetStore) ScanPage(ctx context.Context, key string, cursor uint64, count int64) ([]string, uint64, error) {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return []string{}, 0, nil
	}
	if count <= 0 {
		count = 100
	}
	values, nextCursor, err := s.client.SScan(ctx, redisKey, cursor, "", count).Result()
	if err != nil {
		return nil, 0, err
	}
	return cleanMembers(values), nextCursor, nil
}

func (s *StringSetStore) redisKeyAndMembers(key string, members []string) (string, []string, bool) {
	redisKey, ok := s.redisKey(key)
	if !ok {
		return "", nil, false
	}
	clean := cleanMembers(members)
	return redisKey, clean, len(clean) > 0
}

func (s *StringSetStore) redisKey(key string) (string, bool) {
	if s == nil || s.client == nil {
		return "", false
	}
	return s.keyspace.Key(key)
}

func cleanMembers(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func stringArgs(values []string) []interface{} {
	out := make([]interface{}, 0, len(values))
	for _, value := range values {
		out = append(out, value)
	}
	return out
}
