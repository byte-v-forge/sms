package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func (s *OrderServer) MarkMessageSent(ctx context.Context, request *smsv1.MarkMessageSentRequest) (*smsv1.MarkMessageSentResponse, error) {
	order, err := s.service.MarkMessageSent(ctx, request.GetOrderId(), request.GetRequestId())
	if err != nil {
		return &smsv1.MarkMessageSentResponse{Order: toProtoOrder(order), Error: toProtoError(err)}, nil
	}
	return &smsv1.MarkMessageSentResponse{Order: toProtoOrder(order)}, nil
}

func (s *OrderServer) RequestAdditionalCode(ctx context.Context, request *smsv1.RequestAdditionalCodeRequest) (*smsv1.RequestAdditionalCodeResponse, error) {
	order, err := s.service.RequestAdditionalCode(ctx, request.GetOrderId(), request.GetRequestId())
	if err != nil {
		return &smsv1.RequestAdditionalCodeResponse{Order: toProtoOrder(order), Error: toProtoError(err)}, nil
	}
	return &smsv1.RequestAdditionalCodeResponse{Order: toProtoOrder(order)}, nil
}

func (s *OrderServer) CompleteOrder(ctx context.Context, request *smsv1.CompleteOrderRequest) (*smsv1.CompleteOrderResponse, error) {
	order, err := s.service.CompleteOrder(ctx, request.GetOrderId(), request.GetRequestId())
	if err != nil {
		return &smsv1.CompleteOrderResponse{Order: toProtoOrder(order), Error: toProtoError(err)}, nil
	}
	return &smsv1.CompleteOrderResponse{Order: toProtoOrder(order)}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, request *smsv1.CancelOrderRequest) (*smsv1.CancelOrderResponse, error) {
	order, err := s.service.CancelOrder(ctx, request.GetOrderId(), request.GetRequestId())
	if err != nil {
		return &smsv1.CancelOrderResponse{Order: toProtoOrder(order), Error: toProtoError(err)}, nil
	}
	return &smsv1.CancelOrderResponse{Order: toProtoOrder(order)}, nil
}
