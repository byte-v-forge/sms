package redisx

import "time"

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
