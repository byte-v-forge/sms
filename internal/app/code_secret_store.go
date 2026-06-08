package app

import "github.com/byte-v-forge/sms/internal/core"

type SMSCodeSecretStore struct {
	store TTLStringStore
	clock core.Clock
}

func NewSMSCodeSecretStore(store TTLStringStore, clock core.Clock) *SMSCodeSecretStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &SMSCodeSecretStore{store: store, clock: clock}
}
