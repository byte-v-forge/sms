package herosms

import (
	"context"
	"net/url"
)

type serviceMetadata struct {
	Service string `json:"service"`
	Name    string `json:"name"`
}

type servicesResponse struct {
	Data []serviceMetadata `json:"data"`
}

func (c *Client) ListServices(ctx context.Context) ([]serviceMetadata, error) {
	params := url.Values{}
	params.Set("page", "1")
	params.Set("size", "1000")
	var response servicesResponse
	if err := c.getOpenAPIJSON(ctx, "/left-menu/services", params, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}
