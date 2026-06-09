package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *CatalogService) catalogProviderConfigs(ctx context.Context, providerKeys []string) ([]*smsinternalv1.SmsProviderConfig, error) {
	providerKeys = sortedProviderFilterKeys(normalizedProviderFilter(providerKeys))
	configs, err := s.configs.ListProviderConfigs(ctx, false, singleProviderKey(providerKeys))
	if err != nil {
		return nil, err
	}
	return filteredCatalogProviderConfigs(configs, normalizedProviderFilter(providerKeys)), nil
}
