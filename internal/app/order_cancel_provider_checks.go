package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func validateCancelableOrder(order core.Order) error {
	if order.Status.IsFinal() {
		return core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	return nil
}

func orderHasNoProviderLease(order core.Order) bool {
	return !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == ""
}
