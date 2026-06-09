package app

import (
	"context"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) validateCodeSecretResolveRequest(ctx context.Context, orderID string, ref *commonv1.SecretRef) error {
	if orderID == "" {
		return core.NewError(core.CodeValidationFailed, "order_id is required", false)
	}
	if s == nil || s.store == nil {
		return core.NewError(core.CodeInternal, "sms order store is not configured", true)
	}
	if s.codeSecrets == nil {
		return core.NewError(core.CodeInternal, "sms code secret store is not configured", true)
	}
	if err := validateSMSCodeSecretRef(ref); err != nil {
		return err
	}
	return s.validateCodeSecretAttached(ctx, orderID, ref.GetSecretId())
}

func (s *OrderService) validateCodeSecretAttached(ctx context.Context, orderID string, secretID string) error {
	if _, err := s.GetOrder(ctx, orderID); err != nil {
		return err
	}
	exists, err := s.store.CodeSecretExists(ctx, orderID, secretID)
	if err != nil {
		return core.NewError(core.CodeInternal, "sms code secret attachment check failed", true)
	}
	if !exists {
		return core.NewError(core.CodeValidationFailed, "sms code secret ref is not attached to order", false)
	}
	return nil
}
