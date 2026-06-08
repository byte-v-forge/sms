package redisx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/timex"
	"github.com/redis/go-redis/v9"
)

type BestEffortLocker struct {
	client   redis.Cmdable
	keyspace Keyspace
	ttl      time.Duration
	retry    time.Duration
}

type Lock struct {
	client redis.Cmdable
	key    string
	token  string
}

func NewBestEffortLocker(client redis.Cmdable, prefix string, ttl time.Duration, retry time.Duration) *BestEffortLocker {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if retry <= 0 {
		retry = 100 * time.Millisecond
	}
	return &BestEffortLocker{
		client:   client,
		keyspace: NewKeyspace(prefix),
		ttl:      ttl,
		retry:    retry,
	}
}

func (l *BestEffortLocker) Lock(ctx context.Context, key string) (*Lock, error) {
	redisKey, ok := l.redisKey(key)
	if !ok {
		return nil, fmt.Errorf("redis lock key is required")
	}
	token, err := lockToken()
	if err != nil {
		return nil, err
	}
	for {
		locked, err := l.client.SetNX(ctx, redisKey, token, l.ttl).Result()
		if err != nil {
			return nil, err
		}
		if locked {
			return &Lock{client: l.client, key: redisKey, token: token}, nil
		}
		if err := timex.Sleep(ctx, l.retry); err != nil {
			return nil, err
		}
	}
}

func (l *Lock) Unlock(ctx context.Context) error {
	if l == nil || l.client == nil || l.key == "" || l.token == "" {
		return nil
	}
	return unlockScript.Run(ctx, l.client, []string{l.key}, l.token).Err()
}

func (l *BestEffortLocker) redisKey(key string) (string, bool) {
	if l == nil || l.client == nil {
		return "", false
	}
	return l.keyspace.Key(key)
}

func lockToken() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw[:]), nil
}

var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`)
