package grpcadapter

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminServer) UpsertProviderConfig(ctx context.Context, request *smsinternalv1.UpsertProviderConfigRequest) (*smsinternalv1.UpsertProviderConfigResponse, error) {
	config, err := s.service.UpsertProviderConfig(ctx, request.GetConfig())
	if err != nil {
		return &smsinternalv1.UpsertProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.UpsertProviderConfigResponse{Config: config}, nil
}

func (s *ProviderAdminServer) DeleteProviderConfig(ctx context.Context, request *smsinternalv1.DeleteProviderConfigRequest) (*smsinternalv1.DeleteProviderConfigResponse, error) {
	if err := s.service.DeleteProviderConfig(ctx, request.GetProviderKey()); err != nil {
		return &smsinternalv1.DeleteProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.DeleteProviderConfigResponse{}, nil
}
