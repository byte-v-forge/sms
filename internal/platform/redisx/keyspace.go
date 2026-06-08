package redisx

import "strings"

type Keyspace struct {
	prefix string
}

func NewKeyspace(prefix string) Keyspace {
	return Keyspace{prefix: strings.Trim(strings.TrimSpace(prefix), ":")}
}

func (k Keyspace) Key(key string) (string, bool) {
	_, redisKey, ok := k.CleanKey(key)
	return redisKey, ok
}

func (k Keyspace) CleanKey(key string) (string, string, bool) {
	cleanKey := strings.TrimSpace(key)
	if cleanKey == "" {
		return "", "", false
	}
	if k.prefix == "" {
		return cleanKey, cleanKey, true
	}
	return cleanKey, k.prefix + ":" + cleanKey, true
}
