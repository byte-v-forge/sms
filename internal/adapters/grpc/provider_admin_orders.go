package grpcadapter

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminServer) ListOrders(ctx context.Context, request *smsinternalv1.ListOrdersRequest) (*smsinternalv1.ListOrdersResponse, error) {
	orders, err := s.service.ListOrders(ctx, request.GetIncludeFinal(), int(request.GetLimit()))
	if err != nil {
		return &smsinternalv1.ListOrdersResponse{Error: toProviderError(err)}, nil
	}
	out := make([]*smsinternalv1.SmsOrderView, 0, len(orders))
	for _, order := range orders {
		out = append(out, toOrderView(order))
	}
	return &smsinternalv1.ListOrdersResponse{Orders: out}, nil
}

func (s *ProviderAdminServer) ListOrderCodes(ctx context.Context, request *smsinternalv1.ListOrderCodesRequest) (*smsinternalv1.ListOrderCodesResponse, error) {
	codes, err := s.service.ListOrderCodes(ctx, request.GetOrderIds(), int(request.GetLimitPerOrder()))
	if err != nil {
		return &smsinternalv1.ListOrderCodesResponse{Error: toProviderError(err)}, nil
	}
	out := make([]*smsinternalv1.SmsOrderCodeView, 0, len(codes))
	for _, code := range codes {
		out = append(out, toOrderCodeView(code))
	}
	return &smsinternalv1.ListOrderCodesResponse{Codes: out}, nil
}

func (s *ProviderAdminServer) CancelOrder(ctx context.Context, request *smsinternalv1.CancelProviderOrderRequest) (*smsinternalv1.CancelProviderOrderResponse, error) {
	order, err := s.service.CancelOrder(ctx, request.GetOrderId(), request.GetRequestId())
	if err != nil {
		return &smsinternalv1.CancelProviderOrderResponse{Order: toOrderView(order), Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.CancelProviderOrderResponse{Order: toOrderView(order)}, nil
}
