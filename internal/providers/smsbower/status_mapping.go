package smsbower

import (
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func parseStatus(result string) (core.ProviderCodeResult, error) {
	switch {
	case strings.HasPrefix(result, "STATUS_OK:"):
		return core.ProviderCodeResult{
			Status:     core.StatusCodeReceived,
			Code:       strings.Trim(strings.TrimSpace(strings.TrimPrefix(result, "STATUS_OK:")), "'\""),
			ReceivedAt: time.Now().UTC(),
		}, nil
	case result == "STATUS_WAIT_CODE":
		return core.ProviderCodeResult{Status: core.StatusPendingCode}, nil
	case strings.HasPrefix(result, "STATUS_WAIT_RETRY:"):
		return core.ProviderCodeResult{
			Status: core.StatusAdditionalCodeRequested,
			Code:   strings.TrimSpace(strings.TrimPrefix(result, "STATUS_WAIT_RETRY:")),
		}, nil
	case result == "STATUS_CANCEL":
		return core.ProviderCodeResult{Status: core.StatusCanceled}, nil
	default:
		return core.ProviderCodeResult{}, handlerapi.MapTextError(result)
	}
}

func statusForAction(action core.ProviderAction) (status string, expected string, err error) {
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
		return "", "", core.NewError(core.CodeUnsupportedOperation, "unsupported smsbower status action", false)
	}
}

func isProviderTextError(result string) bool {
	return !strings.HasPrefix(strings.TrimSpace(result), "{")
}
