package handlerapi

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func ActivationStatusForAction(providerName string, action core.ProviderAction) (status string, expected string, err error) {
	switch action {
	case core.ActionMarkMessageSent:
		return "1", "ACCESS_READY", nil
	case core.ActionRequestAdditional:
		return "3", "ACCESS_RETRY_GET", nil
	case core.ActionCompleteOrder:
		return "6", "ACCESS_ACTIVATION", nil
	case core.ActionCancelOrder:
		return "8", "ACCESS_CANCEL", nil
	default:
		providerName = strings.TrimSpace(providerName)
		if providerName == "" {
			providerName = "sms provider"
		}
		return "", "", core.NewError(core.CodeUnsupportedOperation, "unsupported "+providerName+" status action", false)
	}
}
