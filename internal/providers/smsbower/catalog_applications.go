package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func (c *Client) ListApplications(ctx context.Context) ([]ApplicationOffer, error) {
	payload, err := c.getApplicationsPayload(ctx)
	if err != nil {
		return nil, err
	}
	if payload.Status != "" && payload.Status != "success" {
		return nil, handlerapi.MapTextError(payload.Status)
	}
	return decodeApplicationOffers(payload.Services)
}
