package grpcadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

type CatalogServer struct {
	smsv1.UnimplementedSmsCatalogServiceServer
	service *app.CatalogService
}

func NewCatalogServer(service *app.CatalogService) *CatalogServer {
	return &CatalogServer{service: service}
}

func (s *CatalogServer) ListSmsProviders(ctx context.Context, _ *smsv1.ListSmsProvidersRequest) (*smsv1.ListSmsProvidersResponse, error) {
	providers, err := s.service.ListProviders(ctx)
	if err != nil {
		return &smsv1.ListSmsProvidersResponse{Error: toProtoError(err)}, nil
	}
	out := make([]*smsv1.SmsProviderInfo, 0, len(providers))
	for _, provider := range providers {
		capabilities := provider.GetCapabilities()
		out = append(out, &smsv1.SmsProviderInfo{
			ProviderKey:             provider.GetProviderKey(),
			DisplayName:             provider.GetDisplayName(),
			SupportsBalance:         capabilities.GetSupportsBalance(),
			SupportsCatalog:         capabilities.GetSupportsCatalog(),
			SupportsAdditionalCode:  capabilities.GetSupportsAdditionalCode(),
			RequiresMarkMessageSent: capabilities.GetRequiresMarkMessageSent(),
		})
	}
	return &smsv1.ListSmsProvidersResponse{Providers: out}, nil
}
