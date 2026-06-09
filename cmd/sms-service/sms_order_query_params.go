package main

import (
	"net/http"
	"net/url"
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/platform/httpx"
)

func smsListOrdersRequest(r *http.Request) *smsinternalv1.ListOrdersRequest {
	return &smsinternalv1.ListOrdersRequest{
		IncludeFinal: httpx.QueryBool(r, "include_final", false),
		Limit:        int32(httpx.QueryInt(r, "limit", 100)),
	}
}

func smsPriceOffersRequest(r *http.Request) *smsv1.ListSmsPriceOffersRequest {
	query := r.URL.Query()
	return &smsv1.ListSmsPriceOffersRequest{
		ApplicationKey:     strings.TrimSpace(query.Get("application_key")),
		CountryIso2:        strings.ToUpper(strings.TrimSpace(query.Get("country_iso2"))),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(query.Get("country_calling_code")), "+"),
		ProviderKeys:       smsProviderKeysFromQuery(query),
	}
}

func smsListOrderCodesRequest(r *http.Request) *smsinternalv1.ListOrderCodesRequest {
	return &smsinternalv1.ListOrderCodesRequest{
		OrderIds:      orderIDsFromQuery(r),
		LimitPerOrder: int32(httpx.QueryInt(r, "limit_per_order", 10)),
	}
}

func smsProviderKeysFromQuery(query url.Values) []string {
	return trimmedQueryValues(query["provider_key"])
}

func orderIDsFromQuery(r *http.Request) []string {
	query := r.URL.Query()
	values := query["order_id"]
	if csv := strings.TrimSpace(query.Get("order_ids")); csv != "" {
		values = append(values, strings.Split(csv, ",")...)
	}
	return uniqueTrimmedValues(values)
}

func trimmedQueryValues(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			out = append(out, value)
		}
	}
	return out
}

func uniqueTrimmedValues(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		id := strings.TrimSpace(value)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
