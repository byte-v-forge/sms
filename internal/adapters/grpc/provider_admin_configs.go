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
