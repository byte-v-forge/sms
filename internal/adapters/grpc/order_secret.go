package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func (s *OrderServer) ResolveSmsCodeSecret(ctx context.Context, request *smsv1.ResolveSmsCodeSecretRequest) (*smsv1.ResolveSmsCodeSecretResponse, error) {
	value, err := s.service.ResolveCodeSecret(ctx, request.GetOrderId(), request.GetSecretRef())
	if err != nil {
		return &smsv1.ResolveSmsCodeSecretResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.ResolveSmsCodeSecretResponse{CodeValue: value}, nil
}
