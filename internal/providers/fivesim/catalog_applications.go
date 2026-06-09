package fivesim

import (
	"context"
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/searchx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func (c *Client) ListCatalogApplications(ctx context.Context, query core.CatalogApplicationQuery) ([]core.CatalogApplication, error) {
	var raw map[string]struct {
		Category string          `json:"Category"`
		Name     string          `json:"Name"`
		Quantity json.RawMessage `json:"Qty"`
		Price    json.RawMessage `json:"Price"`
	}
	if err := c.getJSON(ctx, "/v1/guest/products/any/any", nil, false, &raw); err != nil {
		return nil, err
	}
	applications := make([]core.CatalogApplication, 0, len(raw))
	for product, item := range raw {
		application := core.CatalogApplication{
			ApplicationKey: product,
			DisplayName:    stringx.FirstNonEmpty(item.Name, fiveSimApplicationName(product)),
			Aliases:        []string{product, item.Name},
		}
		if !catalogApplicationMatches(application, query.SearchText) {
			continue
		}
		applications = append(applications, application)
	}
	return applications, nil
}

func catalogApplicationMatches(application core.CatalogApplication, searchText string) bool {
	if searchx.Token(searchText) == "" {
		return true
	}
	for _, value := range append([]string{application.ApplicationKey, application.DisplayName}, application.Aliases...) {
		if searchx.ContainsToken(value, searchText) {
			return true
		}
	}
	return false
}
