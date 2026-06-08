package main

import (
	"errors"
	"net/http"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *dashboardServer) handleSMSOrder(w http.ResponseWriter, r *http.Request) {
	id, action, ok := splitSMSPath(r.URL.Path, "/orders/")
	if !ok {
		writeError(w, http.StatusNotFound, errors.New("sms order action not found"))
		return
	}
	if !requireHTTPMethod(w, r, http.MethodPost) {
		return
	}
	switch action {
	case "mark-message-sent":
		s.handleSMSOrderMarkMessageSent(w, r, id)
	case "additional-code", "request-additional-code":
		s.handleSMSOrderRequestAdditionalCode(w, r, id)
	case "complete":
		s.handleSMSOrderComplete(w, r, id)
	case "cancel":
		s.handleSMSOrderCancel(w, r, id)
	default:
		writeError(w, http.StatusNotFound, errors.New("sms order action not found"))
	}
}

func (s *dashboardServer) handleSMSOrderMarkMessageSent(w http.ResponseWriter, r *http.Request, orderID string) {
	var req smsv1.MarkMessageSentRequest
	if err := readProtoJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.OrderId = orderID
	resp, err := s.smsOrderClient.MarkMessageSent(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSOrderRequestAdditionalCode(w http.ResponseWriter, r *http.Request, orderID string) {
	var req smsv1.RequestAdditionalCodeRequest
	if err := readProtoJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.OrderId = orderID
	resp, err := s.smsOrderClient.RequestAdditionalCode(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSOrderComplete(w http.ResponseWriter, r *http.Request, orderID string) {
	var req smsv1.CompleteOrderRequest
	if err := readProtoJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.OrderId = orderID
	resp, err := s.smsOrderClient.CompleteOrder(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSOrderCancel(w http.ResponseWriter, r *http.Request, orderID string) {
	var req smsinternalv1.CancelProviderOrderRequest
	if err := readProtoJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.OrderId = orderID
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
