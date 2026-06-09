package fivesim

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func orderStatusResult(status string) core.ProviderCodeResult {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "CANCELED":
		return core.ProviderCodeResult{Status: core.StatusCanceled}
	case "TIMEOUT":
		return core.ProviderCodeResult{Status: core.StatusExpired}
	case "FINISHED":
		return core.ProviderCodeResult{Status: core.StatusCompleted}
	case "BANNED":
		return core.ProviderCodeResult{Status: core.StatusFailed}
	default:
		return core.ProviderCodeResult{Status: core.StatusPendingCode}
	}
}
