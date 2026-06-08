package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func (s *OrderServer) GetOrder(ctx context.Context, request *smsv1.GetOrderRequest) (*smsv1.GetOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, request.GetOrderId())
	if err != nil {
		return &smsv1.GetOrderResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.GetOrderResponse{Order: toProtoOrder(order)}, nil
}
