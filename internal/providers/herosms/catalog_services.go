package herosms

import (
	"context"
	"net/url"
	"strconv"

	"github.com/byte-v-forge/sms/internal/core"
)

const (
	heroSMSCatalogPageSize = 25
	heroSMSMaxCatalogPages = 200
)

type serviceMetadata struct {
	Service string `json:"service"`
	Name    string `json:"name"`
}

type servicesResponse struct {
	Data    []serviceMetadata `json:"data"`
	HasMore bool              `json:"hasMore"`
}

func (c *Client) ListServices(ctx context.Context) ([]serviceMetadata, error) {
	services := make([]serviceMetadata, 0)
	for page := 1; page <= heroSMSMaxCatalogPages; page++ {
		response, err := c.listServicesPage(ctx, page)
		if err != nil {
			return nil, err
		}
		services = append(services, response.Data...)
		if !response.HasMore || len(response.Data) == 0 {
			return services, nil
		}
	}
	return nil, core.NewError(core.CodeSupplyUnavailable, "hero sms service catalog exceeded page limit", true)
}

func (c *Client) listServicesPage(ctx context.Context, page int) (servicesResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(heroSMSCatalogPageSize))
	var response servicesResponse
	if err := c.getOpenAPIJSON(ctx, "/left-menu/services", params, &response); err != nil {
		return servicesResponse{}, err
	}
	return response, nil
}
