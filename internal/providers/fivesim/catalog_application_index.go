package fivesim

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) catalogApplicationNames(ctx context.Context) map[string]string {
	applications, err := c.ListCatalogApplications(ctx)
	if err != nil {
		return nil
	}
	return fiveSimApplicationNameIndex(applications)
}

func fiveSimApplicationNameIndex(applications []core.CatalogApplication) map[string]string {
	names := make(map[string]string, len(applications))
	for _, app := range applications {
		if app.ApplicationKey != "" && app.DisplayName != "" {
			names[app.ApplicationKey] = app.DisplayName
		}
	}
	return names
}
