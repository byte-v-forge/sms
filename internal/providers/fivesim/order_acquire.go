package fivesim

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	country := strings.TrimSpace(request.Route.ProviderCountryID)
	if country == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "5sim country is required", false)
	}
	product := strings.TrimSpace(request.Route.UpstreamServiceKey)
	if product == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "5sim product is required", false)
	}
	operator, err := operatorForRoute(request.Route)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	path := fmt.Sprintf("/v1/user/buy/activation/%s/%s/%s", url.PathEscape(country), url.PathEscape(operator), url.PathEscape(product))

	var payload order
	if err := c.getJSON(ctx, path, nil, true, &payload); err != nil {
		return core.ProviderOrder{}, err
	}
	orderID := jsonx.Scalar(payload.ID)
	if orderID == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "missing 5sim order id", false)
	}
	e164, national := phone.Normalize(payload.Phone, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID: orderID,
		PhoneNumber: core.PhoneNumber{
			E164:               e164,
			NationalNumber:     national,
			CountryISO2:        request.Target.CountryISO2,
			CountryCallingCode: request.Target.CountryCallingCode,
		},
		Price:      core.Money{CurrencyCode: c.currencyCode, AmountDecimal: jsonx.Scalar(payload.Price)},
		AcquiredAt: parseTime(payload.CreatedAt),
		ExpiresAt:  parseTime(payload.Expires),
	}, nil
}
