package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"golang.org/x/sync/errgroup"
)

func (s *CatalogService) ListApplications(ctx context.Context, query core.CatalogApplicationQuery) ([]core.CatalogApplication, error) {
	providerKeys := sortedProviderFilterKeys(normalizedProviderFilter(query.ProviderKeys))
	configs, err := s.catalogProviderConfigs(ctx, providerKeys)
	if err != nil {
		return nil, err
	}
	results := make([]catalogProviderApplicationsResult, len(configs))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(catalogProviderConcurrency)
	for index, config := range configs {
		index := index
		config := config
		group.Go(func() error {
			results[index].applications, results[index].err = s.listProviderCatalogApplications(groupCtx, config)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}
	return collectCatalogProviderApplications(results)
}
