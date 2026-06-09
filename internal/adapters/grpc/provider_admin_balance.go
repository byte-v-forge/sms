package grpcadapter

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminServer) GetProviderBalance(ctx context.Context, request *smsinternalv1.GetProviderBalanceRequest) (*smsinternalv1.GetProviderBalanceResponse, error) {
	balance, err := s.service.GetProviderBalance(ctx, request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.GetProviderBalanceResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetProviderBalanceResponse{Balance: toProtoMoney(balance)}, nil
}
