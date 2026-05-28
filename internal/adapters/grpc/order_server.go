package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

type OrderServer struct {
	smsv1.UnimplementedSmsOrderServiceServer
	service *app.OrderService
}

func NewOrderServer(service *app.OrderService) *OrderServer {
	return &OrderServer{service: service}
}

func (s *OrderServer) AcquireNumber(ctx context.Context, request *smsv1.AcquireNumberRequest) (*smsv1.AcquireNumberResponse, error) {
	order, err := s.service.AcquireNumber(ctx, core.AcquireNumberCommand{
		RequestID:     request.GetRequestId(),
		AcquireParams: app.RouteFromPublicAcquireParams(request.GetAcquireParams()),
		LeaseDuration: protoDuration(request.GetLeaseDuration()),
	})
	if err != nil {
		return &smsv1.AcquireNumberResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.AcquireNumberResponse{Order: toProtoOrder(order)}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, request *smsv1.GetOrderRequest) (*smsv1.GetOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, request.GetOrderId())
	if err != nil {
		return &smsv1.GetOrderResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.GetOrderResponse{Order: toProtoOrder(order)}, nil
}

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
