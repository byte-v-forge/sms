package main

import (
	"errors"
	"net/http"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/httpx"
)

func (s *dashboardServer) handleSMSOrderAcquire(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodPost) {
		return
	}
	if s.smsOrderClient == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("sms order service is not configured"))
		return
	}
	var req smsv1.AcquireNumberRequest
	if err := readProtoJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if smsAcquireParamsNeedRecommendation(req.GetAcquireParams()) {
		params, smsErr := s.recommendSMSAcquireParams(r.Context(), req.GetAcquireParams())
		if smsErr != nil {
			writeProtoJSON(w, http.StatusOK, &smsv1.AcquireNumberResponse{Error: smsErr})
			return
		}
		req.AcquireParams = params
	}
	resp, err := s.smsOrderClient.AcquireNumber(r.Context(), &req)
	if err != nil {
		writeProtoJSON(w, http.StatusOK, &smsv1.AcquireNumberResponse{Error: app.PublicError(err)})
		return
	}
	if resp.GetError() == nil && resp.GetOrder().GetOrderId() != "" {
		resp = s.waitSMSOrderAcquired(r.Context(), resp, time.Duration(httpx.QueryInt(r, "wait_seconds", 60))*time.Second)
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
