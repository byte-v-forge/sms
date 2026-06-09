package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *MemoryOrderStore) CodeSecretExists(ctx context.Context, orderID string, secretID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return memoryCodeSecretExists(s.codes[orderID], secretID), nil
}

func memoryCodeSecretExists(codes []core.OrderCode, secretID string) bool {
	for _, code := range codes {
		if code.Code.SecretRef.GetSecretId() == secretID {
			return true
		}
	}
	return false
}
