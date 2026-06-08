package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) deferCancelRetry(ctx context.Context, order core.Order, requestID string, retryAt time.Time) (core.Order, error) {
	if retryAt.IsZero() {
		retryAt = s.clock.Now().Add(5 * time.Second)
	}
	if !retryAt.After(s.clock.Now()) {
		retryAt = s.clock.Now().Add(5 * time.Second)
	}
	order = markCancelRequest(order, requestID, retryAt, s.clock.Now())
	if err := s.updateOrder(ctx, order); err != nil {
		return core.Order{}, err
	}
	return order, &CancelRetryError{RetryAt: retryAt}
}
