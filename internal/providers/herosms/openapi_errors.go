package herosms

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func mapHeroSMSOpenAPIError(statusCode int, text string) error {
	var payload struct {
		Title   string `json:"title"`
		Details string `json:"details"`
	}
	if err := json.Unmarshal([]byte(text), &payload); err == nil && strings.TrimSpace(payload.Title) != "" {
		return handlerapi.MapTextError(text)
	}
	if text != "" {
		return handlerapi.MapTextError(text)
	}
	if statusCode >= 500 {
		return core.NewError(core.CodeSupplyUnavailable, fmt.Sprintf("hero sms openapi http status %d", statusCode), true)
	}
	return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("hero sms openapi http status %d", statusCode), false)
}
