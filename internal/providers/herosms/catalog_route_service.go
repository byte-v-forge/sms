package herosms

import (
	"context"
	"strings"
)

type routeService struct {
	Key      string
	NameByID map[string]string
}

func (c *Client) routeService(ctx context.Context, applicationKey string) (routeService, error) {
	applicationKey = strings.TrimSpace(applicationKey)
	if applicationKey == "" {
		return routeService{}, nil
	}
	services, err := c.SearchServices(ctx, applicationKey)
	if err != nil {
		return routeService{}, err
	}
	service := heroSMSServiceForQuery(applicationKey, services)
	if service == "" {
		service = normalizeHeroSMSServiceKey(applicationKey)
	}
	return routeService{Key: service, NameByID: heroSMSServiceNameIndex(services)}, nil
}
