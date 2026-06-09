package grpcadapter

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminServer) GetProviderConfig(ctx context.Context, request *smsinternalv1.GetProviderConfigRequest) (*smsinternalv1.GetProviderConfigResponse, error) {
	config, err := s.service.GetProviderConfig(ctx, request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.GetProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetProviderConfigResponse{Config: config}, nil
}

func (s *ProviderAdminServer) ListProviderConfigs(ctx context.Context, request *smsinternalv1.ListProviderConfigsRequest) (*smsinternalv1.ListProviderConfigsResponse, error) {
	configs, err := s.service.ListProviderConfigs(ctx, request.GetIncludeDisabled(), request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.ListProviderConfigsResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.ListProviderConfigsResponse{Configs: configs}, nil
}
