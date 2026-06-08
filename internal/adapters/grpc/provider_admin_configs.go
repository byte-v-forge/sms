package grpcadapter

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminServer) ListProviderPlugins(ctx context.Context, request *smsinternalv1.ListProviderPluginsRequest) (*smsinternalv1.ListProviderPluginsResponse, error) {
	plugins, err := s.service.ListProviderPlugins(ctx)
	if err != nil {
		return &smsinternalv1.ListProviderPluginsResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.ListProviderPluginsResponse{Plugins: plugins}, nil
}

func (s *ProviderAdminServer) UpsertProviderConfig(ctx context.Context, request *smsinternalv1.UpsertProviderConfigRequest) (*smsinternalv1.UpsertProviderConfigResponse, error) {
	config, err := s.service.UpsertProviderConfig(ctx, request.GetConfig())
	if err != nil {
		return &smsinternalv1.UpsertProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.UpsertProviderConfigResponse{Config: config}, nil
}

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

func (s *ProviderAdminServer) DeleteProviderConfig(ctx context.Context, request *smsinternalv1.DeleteProviderConfigRequest) (*smsinternalv1.DeleteProviderConfigResponse, error) {
	if err := s.service.DeleteProviderConfig(ctx, request.GetProviderKey()); err != nil {
		return &smsinternalv1.DeleteProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.DeleteProviderConfigResponse{}, nil
}

func (s *ProviderAdminServer) GetProviderBalance(ctx context.Context, request *smsinternalv1.GetProviderBalanceRequest) (*smsinternalv1.GetProviderBalanceResponse, error) {
	balance, err := s.service.GetProviderBalance(ctx, request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.GetProviderBalanceResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetProviderBalanceResponse{Balance: toProtoMoney(balance)}, nil
}
