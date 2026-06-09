package herosms

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/byte-v-forge/sms/internal/core"
)

const heroSMSActivationTypeSMS = 0

type activationPurchaseRequest struct {
	Owner          int             `json:"owner"`
	FixedPrice     bool            `json:"fixedPrice"`
	ActivationType int             `json:"activationType"`
	Service        string          `json:"service"`
	Country        int             `json:"country"`
	Operator       string          `json:"operator"`
	Cost           json.RawMessage `json:"cost,omitempty"`
	Amount         int             `json:"amount"`
}

func (c *Client) buyActivation(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	payload, err := c.postOpenAPIJSON(ctx, "/activations", activationPurchaseFromRoute(request.Route))
	if err != nil {
		return core.ProviderOrder{}, err
	}
	order, err := decodeActivationPurchase(payload)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	return heroSMSProviderOrderFromPayload(order, request)
}

func activationPurchaseFromRoute(route core.Route) activationPurchaseRequest {
	cost := heroSMSCost(route.MaxPrice.AmountDecimal)
	return activationPurchaseRequest{
		Owner:          6,
		FixedPrice:     cost != nil,
		ActivationType: heroSMSActivationTypeSMS,
		Service:        route.UpstreamServiceKey,
		Country:        heroSMSCountryID(route.ProviderCountryID),
		Operator:       heroSMSOperator(route.UpstreamProviderID),
		Cost:           cost,
		Amount:         1,
	}
}

func heroSMSCountryID(value string) int {
	id, _ := strconv.Atoi(value)
	return id
}

func heroSMSOperator(value string) string {
	if value != "" {
		return value
	}
	return "any"
}

func heroSMSCost(value string) json.RawMessage {
	if value == "" {
		return nil
	}
	return json.RawMessage(value)
}
