package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) resolveCancelRace(ctx context.Context, order core.Order, provider core.Provider, cancelErr *core.Error) (core.Order, bool, error) {
	if !cancelRaceShouldSync(cancelErr) {
		return order, false, nil
	}
	synced, handled, err := s.syncTerminalOrChargedProviderState(ctx, order, provider)
	if err != nil {
		return order, false, nil
	}
	return synced, handled, nil
}

func cancelRaceShouldSync(cancelErr *core.Error) bool {
	if cancelErr == nil {
		return false
	}
	return !cancelErr.Retryable || cancelErr.Code == core.CodeCancelNotAllowed
}
