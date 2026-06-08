package herosms

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
			Code:       strings.TrimSpace(strings.TrimPrefix(result, "STATUS_OK:")),
			ReceivedAt: time.Now().UTC(),
		}, nil
	case result == "STATUS_WAIT_CODE":
		return core.ProviderCodeResult{Status: core.StatusPendingCode}, nil
	case result == "STATUS_WAIT_RESEND":
		return core.ProviderCodeResult{Status: core.StatusPendingCode}, nil
	case strings.HasPrefix(result, "STATUS_WAIT_RETRY"):
		return core.ProviderCodeResult{
			Status: core.StatusAdditionalCodeRequested,
			Code:   strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(result, "STATUS_WAIT_RETRY"), ":")),
		}, nil
	case result == "STATUS_CANCEL":
		return core.ProviderCodeResult{Status: core.StatusCanceled}, nil
	default:
		return core.ProviderCodeResult{}, handlerapi.MapTextError(result)
	}
}
