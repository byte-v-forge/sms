package herosms

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func (c *Client) ListCatalogApplications(ctx context.Context) ([]core.CatalogApplication, error) {
	serviceKeys, err := c.ListActivationServiceKeys(ctx)
	if err != nil {
		return nil, err
	}
	names := map[string]string{}
	if services, err := c.ListServices(ctx); err == nil {
		names = heroSMSServiceNameIndex(services)
	}
	applications := make([]core.CatalogApplication, 0, len(serviceKeys))
	for _, key := range serviceKeys {
		applications = append(applications, core.CatalogApplication{
			ApplicationKey: key,
			DisplayName:    stringx.FirstNonEmpty(names[key], key),
		})
	}
	return applications, nil
}
