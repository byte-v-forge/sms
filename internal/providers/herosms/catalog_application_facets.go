package herosms

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCatalogApplications(ctx context.Context, query core.CatalogApplicationQuery) ([]core.CatalogApplication, error) {
	if strings.TrimSpace(query.SearchText) == "" {
		return nil, nil
	}
	services, err := c.SearchServices(ctx, query.SearchText)
	if err != nil {
		return nil, err
	}
	applications := make([]core.CatalogApplication, 0, len(services))
	for _, service := range services {
		key := normalizeHeroSMSServiceKey(service.Service)
		if key == "" || key == "full" {
			continue
		}
		applications = append(applications, core.CatalogApplication{
			ApplicationKey: key,
			DisplayName:    service.Name,
			Aliases:        []string{key, service.Service, service.Name},
		})
	}
	return applications, nil
}
