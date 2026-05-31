package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *ProviderAdminService) GetProviderBalance(ctx context.Context, providerKey string) (core.Money, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerKey)
	if err != nil {
		return core.Money{}, err
	}
	if !config.GetEnabled() {
		return core.Money{}, core.NewError(core.CodeValidationFailed, "sms provider config is disabled", false)
	}
	provider, err := providerFromConfig(config, s.timeout, s.defaultHTTPProxy)
	if err != nil {
		return core.Money{}, err
	}
	return provider.GetBalance(ctx)
}

func (s *ProviderAdminService) ListOrders(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if s.orderDB == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms order list is not available", false)
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	return s.orderDB.List(ctx, includeFinal, limit)
}

func (s *ProviderAdminService) ListOrderCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if s.orderDB == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms code history is not available", false)
	}
	ids := uniqueNonEmpty(orderIDs)
	if len(ids) == 0 {
		return []core.OrderCode{}, nil
	}
	if len(ids) > 200 {
		return nil, core.NewError(core.CodeValidationFailed, "too many order_ids", false)
	}
	if limitPerOrder <= 0 {
		limitPerOrder = 10
	}
	if limitPerOrder > 50 {
		limitPerOrder = 50
	}
	return s.orderDB.ListCodes(ctx, ids, limitPerOrder)
}

func (s *ProviderAdminService) CancelOrder(ctx context.Context, orderID string, requestID string) (core.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return core.Order{}, core.NewError(core.CodeValidationFailed, "order_id is required", false)
	}
	if strings.TrimSpace(requestID) == "" {
		requestID = RandomIDGenerator{}.NewID("req_")
	}
	return s.orders.CancelOrder(ctx, orderID, requestID)
}

func uniqueNonEmpty(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}
