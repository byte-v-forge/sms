package herosms

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCatalogApplications(ctx context.Context) ([]core.CatalogApplication, error) {
	services, err := c.ListServices(ctx)
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
