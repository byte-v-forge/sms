package app

import (
	"context"
	"strings"

	commonv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/common-lib/secretref"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) ResolveCodeSecret(ctx context.Context, orderID string, ref *commonv1.SecretRef) (string, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return "", core.NewError(core.CodeValidationFailed, "order_id is required", false)
	}
	if s == nil || s.store == nil {
		return "", core.NewError(core.CodeInternal, "sms order store is not configured", true)
	}
	if s.codeSecrets == nil {
		return "", core.NewError(core.CodeInternal, "sms code secret store is not configured", true)
	}
	if err := validateSMSCodeSecretRef(ref); err != nil {
		return "", err
	}
	if _, err := s.GetOrder(ctx, orderID); err != nil {
		return "", err
	}
	exists, err := s.store.CodeSecretExists(ctx, orderID, ref.GetSecretId())
	if err != nil {
		return "", core.NewError(core.CodeInternal, "sms code secret attachment check failed", true)
	}
	if !exists {
		return "", core.NewError(core.CodeValidationFailed, "sms code secret ref is not attached to order", false)
	}
	value, err := s.codeSecrets.ResolveSecret(ctx, ref)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(value) == "" {
		return "", core.NewError(core.CodeOrderExpired, "sms code secret expired", false)
	}
	return value, nil
}

func validateSMSCodeSecretRef(ref *commonv1.SecretRef) error {
	if err := secretref.Validate(ref); err != nil {
		return core.NewError(core.CodeValidationFailed, "sms code secret ref is invalid", false)
	}
	if ref.GetProvider() != "sms" || ref.GetPurpose() != "sms_otp" {
		return core.NewError(core.CodeValidationFailed, "sms secret ref scope mismatch", false)
	}
	return nil
}
