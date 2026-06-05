package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/common-lib/redisx"
	"github.com/byte-v-forge/common-lib/secretref"
	"github.com/byte-v-forge/sms/internal/core"
)

type SMSCodeSecretStore struct {
	store *redisx.StringStore
	clock core.Clock
}

func NewSMSCodeSecretStore(store *redisx.StringStore, clock core.Clock) *SMSCodeSecretStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &SMSCodeSecretStore{store: store, clock: clock}
}

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

func (s *SMSCodeSecretStore) ResolveSecret(ctx context.Context, ref *commonv1.SecretRef) (string, error) {
	if s == nil || s.store == nil {
		return "", core.NewError(core.CodeInternal, "sms code secret store is not configured", true)
	}
	if err := secretref.Validate(ref); err != nil {
		return "", core.NewError(core.CodeValidationFailed, "sms code secret ref is invalid", false)
	}
	if ref.GetProvider() != "sms" || ref.GetPurpose() != "sms_otp" {
		return "", core.NewError(core.CodeValidationFailed, "sms secret ref scope mismatch", false)
	}
	value, ok, err := s.store.Load(ctx, ref.GetSecretId())
	if err != nil {
		return "", core.NewError(core.CodeInternal, "sms code secret store read failed", true)
	}
	if !ok {
		return "", core.NewError(core.CodeOrderExpired, "sms code secret expired", false)
	}
	return value, nil
}

func smsCodeSecretRef(orderID string, receivedAt time.Time, expiresAt time.Time) *commonv1.SecretRef {
	secretID := secretref.StableID("sms-code", orderID, fmt.Sprintf("%d", receivedAt.UnixNano()))
	return secretref.New("sms", "sms_otp", secretID, expiresAt)
}
