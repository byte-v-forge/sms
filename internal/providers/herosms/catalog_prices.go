package herosms

import (
	"context"
	"net/url"
	"strings"
)

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID string) ([]PriceOffer, error) {
	service := normalizeHeroSMSServiceKey(serviceKey)
	countryID = strings.TrimSpace(countryID)
	if service == "" || countryID == "" || service == "full" {
		return nil, nil
	}
	info, err := c.purchaseInfo(ctx, service, countryID)
	if err != nil {
		return nil, err
	}
	return purchaseInfoPriceOffers(service, countryID, info), nil
}

func (c *Client) purchaseInfo(ctx context.Context, serviceKey, countryID string) (purchaseInfo, error) {
	path := "/left-menu/service/" + url.PathEscape(serviceKey) + "/country/" + url.PathEscape(countryID) + "/offers"
	var response purchaseInfoResponse
	if err := c.getOpenAPIJSON(ctx, path, nil, &response); err != nil {
		return purchaseInfo{}, err
	}
	return purchaseInfoForService(response, serviceKey), nil
}
