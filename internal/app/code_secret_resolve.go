package app

import (
	"context"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
)

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
