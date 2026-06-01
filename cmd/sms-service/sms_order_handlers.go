package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/httpx"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"google.golang.org/protobuf/proto"
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

func (s *dashboardServer) handleSMSOrderAcquire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

func (s *dashboardServer) waitSMSOrderAcquired(ctx context.Context, initial *smsv1.AcquireNumberResponse, timeout time.Duration) *smsv1.AcquireNumberResponse {
	if timeout <= 0 || initial.GetOrder().GetStatus() != smsv1.SmsOrderStatus_SMS_ORDER_STATUS_ACQUIRE_REQUESTED {
		return initial
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	latest := initial.GetOrder()
	for {
		select {
		case <-ctx.Done():
			return &smsv1.AcquireNumberResponse{Order: latest, Error: &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_TIMEOUT, Message: "sms number acquisition timed out", Retryable: true}}
		case <-ticker.C:
		}
		resp, err := s.smsOrderClient.GetOrder(ctx, &smsv1.GetOrderRequest{OrderId: initial.GetOrder().GetOrderId()})
		if err != nil {
			return &smsv1.AcquireNumberResponse{Order: latest, Error: app.PublicError(err)}
		}
		if resp.GetOrder() != nil {
			latest = resp.GetOrder()
		}
		if resp.GetError() != nil {
			return &smsv1.AcquireNumberResponse{Order: latest, Error: resp.GetError()}
		}
		if smsOrderHasPhone(latest) {
			return &smsv1.AcquireNumberResponse{Order: latest}
		}
		if latest.GetStatus() != smsv1.SmsOrderStatus_SMS_ORDER_STATUS_ACQUIRE_REQUESTED {
			if latest.GetLastError() != nil {
				return &smsv1.AcquireNumberResponse{Order: latest, Error: latest.GetLastError()}
			}
			return &smsv1.AcquireNumberResponse{Order: latest, Error: &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_SUPPLY_UNAVAILABLE, Message: "sms number acquisition finished without phone: " + latest.GetStatus().String(), Retryable: true}}
		}
	}
}

func smsOrderHasPhone(order *smsv1.SmsOrder) bool {
	if order == nil || order.GetPhoneNumber() == nil {
		return false
	}
	phone := order.GetPhoneNumber()
	return strings.TrimSpace(phone.GetE164Number()) != "" || strings.TrimSpace(phone.GetNationalNumber()) != ""
}

func smsAcquireParamsNeedRecommendation(params *smsv1.SmsNumberAcquireParams) bool {
	if params == nil {
		return true
	}
	return params.GetProviderParams() == nil
}

func (s *dashboardServer) recommendSMSAcquireParams(ctx context.Context, params *smsv1.SmsNumberAcquireParams) (*smsv1.SmsNumberAcquireParams, *smsv1.SmsError) {
	if s.smsCatalogClient == nil {
		return nil, &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_INTERNAL, Message: "sms catalog service is not configured"}
	}
	target := &smsv1.SmsTarget{
		ApplicationKey:     strings.TrimSpace(params.GetApplicationKey()),
		CountryIso2:        strings.ToUpper(strings.TrimSpace(params.GetCountryIso2())),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(params.GetCountryCallingCode()), "+"),
	}
	resp, err := s.smsCatalogClient.RecommendSmsRoutes(ctx, &smsv1.RecommendSmsRoutesRequest{
		Target: target,
		Policy: &smsv1.SmsRoutePolicy{
			Strategy:          smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE,
			Limit:             1,
			MinAvailableCount: smsMinAvailableCount(params.GetMinAvailableCount()),
			MinPrice:          params.GetMinPrice(),
			MaxPrice:          params.GetMaxPrice(),
			FailurePolicy:     params.GetRouteFailurePolicy(),
		},
	})
	if err != nil {
		return nil, app.PublicError(err)
	}
	if resp.GetError() != nil {
		return nil, resp.GetError()
	}
	if len(resp.GetRecommendations()) == 0 || resp.GetRecommendations()[0].GetOffer().GetAcquireParams() == nil {
		return nil, &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_ROUTE_NOT_FOUND, Message: fmt.Sprintf("sms route not found for %s/%s/%s", target.GetApplicationKey(), target.GetCountryIso2(), target.GetCountryCallingCode()), Retryable: true}
	}
	recommended, ok := proto.Clone(resp.GetRecommendations()[0].GetOffer().GetAcquireParams()).(*smsv1.SmsNumberAcquireParams)
	if !ok || recommended == nil {
		return nil, &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_INTERNAL, Message: "sms route recommendation is invalid"}
	}
	if recommended.RouteFailurePolicy == nil {
		recommended.RouteFailurePolicy = params.GetRouteFailurePolicy()
	}
	return recommended, nil
}

func smsMinAvailableCount(value int32) int32 {
	if value <= 0 {
		return 1
	}
	return value
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
