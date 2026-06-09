package app

import (
	"context"
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) ResolveCodeSecret(ctx context.Context, orderID string, ref *commonv1.SecretRef) (string, error) {
	orderID = strings.TrimSpace(orderID)
	if err := s.validateCodeSecretResolveRequest(ctx, orderID, ref); err != nil {
		return "", err
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
