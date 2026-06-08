package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
)

func (s *OrderService) prepareCodeSecret(ctx context.Context, order core.Order, code core.SMSCode) (core.SMSCode, error) {
	if secretref.Configured(code.SecretRef) {
		return core.SMSCode{ReceivedAt: code.ReceivedAt, SecretRef: code.SecretRef}, nil
	}
	if s == nil || s.codeSecrets == nil {
		return core.SMSCode{}, core.NewError(core.CodeInternal, "sms code secret store is not configured", false)
	}
	return s.codeSecrets.Save(ctx, order, code)
}
