package main

import (
	"errors"
	"net/http"

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
