package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) listProviderCatalogApplications(ctx context.Context, config *smsinternalv1.SmsProviderConfig) ([]core.CatalogApplication, error) {
	if !config.GetEnabled() {
		return nil, nil
	}
	provider, err := providerFromConfig(s.providers, config, s.timeout, s.defaultHTTPProxy)
	if err != nil {
		return nil, err
	}
	catalogProvider, ok := provider.(applicationCatalogProvider)
	if !ok {
		return nil, nil
	}
	return catalogProvider.ListCatalogApplications(ctx)
}
