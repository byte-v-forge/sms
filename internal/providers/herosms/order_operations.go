package herosms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	result, err := c.api.GetNumberV2(ctx, request, handlerapi.GetNumberV2Config{
		ProviderName:    "hero sms",
		CountryLabel:    "provider country id",
		ProviderIDParam: "operator",
		MaxPriceParam:   "maxPrice",
	})
	if err != nil {
		return core.ProviderOrder{}, err
	}
	var payload heroSMSGetNumberV2Response
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		if orderID, rawPhone, ok := parseAccessNumber(result); ok {
			return heroSMSProviderOrder(orderID, rawPhone, heroSMSGetNumberV2Response{}, request), nil
		}
		return core.ProviderOrder{}, handlerapi.MapTextError(result)
	}
	orderID := jsonx.FirstScalar(payload.ActivationID)
	if orderID == "" || strings.TrimSpace(payload.PhoneNumber) == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "bad hero sms getNumberV2 response", false)
	}
	return heroSMSProviderOrder(orderID, payload.PhoneNumber, payload, request), nil
}

func parseAccessNumber(result string) (orderID, rawPhone string, ok bool) {
	parts := strings.SplitN(result, ":", 3)
	if len(parts) != 3 || parts[0] != "ACCESS_NUMBER" {
		return "", "", false
	}
	return parts[1], parts[2], true
}
