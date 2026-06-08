package app

import (
	"context"
	"errors"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *SMSCodeSecretStore) Save(ctx context.Context, order core.Order, code core.SMSCode) (core.SMSCode, error) {
	if s == nil || s.store == nil {
		return core.SMSCode{}, errors.New("sms code secret store is not configured")
	}
	if strings.TrimSpace(code.Value) == "" {
		return core.SMSCode{}, errors.New("sms code value is required")
	}
	receivedAt := code.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = s.clock.Now()
	}
	now := s.clock.Now()
	expiresAt := order.ExpiresAt
	if !expiresAt.After(now) {
		expiresAt = now.Add(s.store.DefaultTTL())
	}
	ref := smsCodeSecretRef(order.ID, receivedAt, expiresAt)
	if ref == nil {
		return core.SMSCode{}, errors.New("sms code secret ref is required")
	}
	ttl := expiresAt.Sub(now)
	if ttl <= 0 {
		return core.SMSCode{}, errors.New("sms code secret ttl is expired")
	}
	if err := s.store.SaveTTL(ctx, ref.GetSecretId(), code.Value, ttl); err != nil {
		return core.SMSCode{}, err
	}
	return core.SMSCode{ReceivedAt: receivedAt, SecretRef: ref}, nil
}
