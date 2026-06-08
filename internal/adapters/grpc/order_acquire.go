package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

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
