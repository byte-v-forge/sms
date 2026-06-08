package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/httpx"
	"google.golang.org/protobuf/proto"
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
			if smsOrderHasPhone(latest) {
				return &smsv1.AcquireNumberResponse{Order: latest}
			}
			return &smsv1.AcquireNumberResponse{Order: latest, Error: app.PublicError(err)}
		}
		if resp.GetOrder() != nil {
			latest = resp.GetOrder()
		}
		if smsOrderHasPhone(latest) {
			return &smsv1.AcquireNumberResponse{Order: latest}
		}
		if resp.GetError() != nil {
			return &smsv1.AcquireNumberResponse{Order: latest, Error: resp.GetError()}
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
	ref := params.GetOfferRef()
	if ref == nil {
		return true
	}
	target := ref.GetTarget()
	routeRef := ref.GetRouteRef()
	return strings.TrimSpace(ref.GetProviderKey()) == "" ||
		target == nil ||
		strings.TrimSpace(target.GetApplicationKey()) == "" ||
		(strings.TrimSpace(target.GetCountryIso2()) == "" && strings.TrimSpace(target.GetCountryCallingCode()) == "") ||
		routeRef == nil ||
		strings.TrimSpace(routeRef.GetUpstreamServiceKey()) == "" ||
		strings.TrimSpace(routeRef.GetProviderCountryId()) == ""
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
	if len(resp.GetRecommendations()) == 0 || resp.GetRecommendations()[0].GetOffer().GetOfferRef() == nil {
		return nil, &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_ROUTE_NOT_FOUND, Message: fmt.Sprintf("sms route not found for %s/%s/%s", target.GetApplicationKey(), target.GetCountryIso2(), target.GetCountryCallingCode()), Retryable: true}
	}
	offerRef, ok := proto.Clone(resp.GetRecommendations()[0].GetOffer().GetOfferRef()).(*smsv1.SmsOfferRef)
	if !ok || offerRef == nil {
		return nil, &smsv1.SmsError{Code: smsv1.SmsErrorCode_SMS_ERROR_CODE_INTERNAL, Message: "sms route recommendation is invalid"}
	}
	recommended := &smsv1.SmsNumberAcquireParams{
		OfferRef:           offerRef,
		ApplicationKey:     target.GetApplicationKey(),
		CountryIso2:        target.GetCountryIso2(),
		CountryCallingCode: target.GetCountryCallingCode(),
		MinAvailableCount:  params.GetMinAvailableCount(),
		MinPrice:           params.GetMinPrice(),
		MaxPrice:           params.GetMaxPrice(),
		RouteFailurePolicy: params.GetRouteFailurePolicy(),
	}
	return recommended, nil
}

func smsMinAvailableCount(value int32) int32 {
	if value <= 0 {
		return 1
	}
	return value
}
