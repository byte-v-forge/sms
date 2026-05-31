package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/byte-v-forge/common-lib/httpx"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *dashboardServer) handleSMSOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

func (s *dashboardServer) handleSMSOrderCodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

func (s *dashboardServer) handleSMSOrder(w http.ResponseWriter, r *http.Request) {
	id, action, ok := splitSMSPath(r.URL.Path, "/orders/")
	if !ok || action != "cancel" {
		writeError(w, http.StatusNotFound, errors.New("sms order action not found"))
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req smsinternalv1.CancelProviderOrderRequest
	if err := readProtoJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.OrderId = id
	resp, err := s.smsAdminClient.CancelOrder(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
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
