package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	catalog, err := c.routeCatalogInputs(ctx, query)
	if err != nil || catalog.empty {
		return nil, err
	}
	priceOffers, err := c.ListPriceOffers(ctx, catalog.serviceKey, catalog.countryID)
	if err != nil {
		return nil, err
	}
	return smsbowerRouteOffers(priceOffers, catalog), nil
}
