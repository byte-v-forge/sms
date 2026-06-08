package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func orderCodeCheckPrecondition(order core.Order) (bool, error) {
	if order.Status == core.StatusCodeReceived {
		return true, nil
	}
	if order.Status == core.StatusAcquireRequested || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return true, nil
	}
	if order.Status.IsFinal() {
		return true, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	return false, nil
}
