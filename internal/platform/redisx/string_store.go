package redisx

import (
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
