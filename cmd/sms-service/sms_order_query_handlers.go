package main

import (
	"errors"
	"net/http"
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/httpx"
)

func (s *dashboardServer) handleSMSOrders(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	resp, err := s.smsAdminClient.ListOrders(r.Context(), &smsinternalv1.ListOrdersRequest{
		IncludeFinal: httpx.QueryBool(r, "include_final", false),
		Limit:        int32(httpx.QueryInt(r, "limit", 100)),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSPriceOffers(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if s.smsCatalogClient == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("sms catalog service is not configured"))
		return
	}
	query := r.URL.Query()
	resp, err := s.smsCatalogClient.ListSmsPriceOffers(r.Context(), &smsv1.ListSmsPriceOffersRequest{
		ApplicationKey:     strings.TrimSpace(query.Get("application_key")),
		CountryIso2:        strings.ToUpper(strings.TrimSpace(query.Get("country_iso2"))),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(query.Get("country_calling_code")), "+"),
		ProviderKeys:       smsProviderKeysFromQuery(query),
	})
	if err != nil {
		writeProtoJSON(w, http.StatusOK, &smsv1.ListSmsPriceOffersResponse{Error: app.PublicError(err)})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSOrderCodes(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	orderIDs := orderIDsFromQuery(r)
	resp, err := s.smsAdminClient.ListOrderCodes(r.Context(), &smsinternalv1.ListOrderCodesRequest{
		OrderIds:      orderIDs,
		LimitPerOrder: int32(httpx.QueryInt(r, "limit_per_order", 10)),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func smsProviderKeysFromQuery(query map[string][]string) []string {
	values := query["provider_key"]
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			out = append(out, value)
		}
	}
	return out
}

func orderIDsFromQuery(r *http.Request) []string {
	values := r.URL.Query()["order_id"]
	if csv := strings.TrimSpace(r.URL.Query().Get("order_ids")); csv != "" {
		values = append(values, strings.Split(csv, ",")...)
	}
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
