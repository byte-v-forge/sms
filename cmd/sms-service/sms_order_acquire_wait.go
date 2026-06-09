package main

import (
	"context"
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

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
