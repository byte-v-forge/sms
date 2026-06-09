package herosms

import (
	"context"
	"net/url"
	"strconv"
	"strings"

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
		response, err := c.listServicesPage(ctx, page, "")
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

func (c *Client) SearchServices(ctx context.Context, query string) ([]serviceMetadata, error) {
	search := heroSMSServiceSearchText(query)
	if search == "" {
		return nil, nil
	}
	response, err := c.listServicesPage(ctx, 1, search)
	if err != nil {
		return nil, err
	}
	if len(response.Data) > 0 {
		return response.Data, nil
	}
	fallback := heroSMSFirstServiceSearchWord(search)
	if fallback == "" || fallback == search {
		return response.Data, nil
	}
	response, err = c.listServicesPage(ctx, 1, fallback)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *Client) listServicesPage(ctx context.Context, page int, search string) (servicesResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(heroSMSCatalogPageSize))
	if search = strings.TrimSpace(search); search != "" {
		params.Set("search", search)
	}
	var response servicesResponse
	if err := c.getOpenAPIJSON(ctx, "/left-menu/services", params, &response); err != nil {
		return servicesResponse{}, err
	}
	return response, nil
}

func heroSMSServiceSearchText(value string) string {
	var out strings.Builder
	for _, r := range strings.TrimSpace(value) {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			out.WriteRune(r)
		default:
			out.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(out.String()), " ")
}

func heroSMSFirstServiceSearchWord(value string) string {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
