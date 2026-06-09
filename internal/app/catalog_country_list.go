package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"golang.org/x/sync/errgroup"
)

func (s *CatalogService) ListCountries(ctx context.Context, query core.CatalogCountryQuery) ([]core.CatalogCountry, error) {
	providerKeys := sortedProviderFilterKeys(normalizedProviderFilter(query.ProviderKeys))
	configs, err := s.catalogProviderConfigs(ctx, providerKeys)
	if err != nil {
		return nil, err
	}
	results := make([]catalogProviderCountriesResult, len(configs))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(catalogProviderConcurrency)
	for index, config := range configs {
		index := index
		config := config
		group.Go(func() error {
			results[index].countries, results[index].err = s.listProviderCatalogCountries(groupCtx, config, routeText(query.ApplicationKey))
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}
	return collectCatalogProviderCountries(results)
}
