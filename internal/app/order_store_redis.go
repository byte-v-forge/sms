package app

import "github.com/byte-v-forge/sms/internal/core"

type RedisOrderStore struct {
	store TTLStringStore
	clock core.Clock
}

func NewRedisOrderStore(store TTLStringStore, clock core.Clock) *RedisOrderStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &RedisOrderStore{store: store, clock: clock}
}
