package smsbower

import (
	"context"
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func (c *Client) ListApplications(ctx context.Context) ([]ApplicationOffer, error) {
	result, err := c.api.Do(ctx, "getServicesList", nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Status   string          `json:"status"`
		Services json.RawMessage `json:"services"`
	}
	if err := decodeJSONObject(result, &payload); err != nil {
		return nil, err
	}
	if payload.Status != "" && payload.Status != "success" {
		return nil, handlerapi.MapTextError(payload.Status)
	}
	offers, err := decodeApplicationOffers(payload.Services)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func decodeApplicationOffers(raw json.RawMessage) ([]ApplicationOffer, error) {
	if len(raw) == 0 {
		return nil, core.NewError(core.CodeUpstreamRejected, "smsbower services list is empty", false)
	}
	var list []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &list); err == nil {
		return applicationOffersFromList(list), nil
	}
	var byCode map[string]struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &byCode); err == nil {
		offers := make([]ApplicationOffer, 0, len(byCode))
		for code, item := range byCode {
			offers = append(offers, applicationOffer(stringx.FirstNonEmpty(item.Code, code), stringx.FirstNonEmpty(item.Name, code)))
		}
		return offers, nil
	}
	var names map[string]string
	if err := json.Unmarshal(raw, &names); err == nil {
		offers := make([]ApplicationOffer, 0, len(names))
		for code, name := range names {
			offers = append(offers, applicationOffer(code, name))
		}
		return offers, nil
	}
	return nil, core.NewError(core.CodeUpstreamRejected, "bad smsbower services list response", false)
}

func applicationOffersFromList(list []struct {
	Code string `json:"code"`
	Name string `json:"name"`
}) []ApplicationOffer {
	offers := make([]ApplicationOffer, 0, len(list))
	for _, service := range list {
		offers = append(offers, applicationOffer(service.Code, service.Name))
	}
	return offers
}

func applicationOffer(code, name string) ApplicationOffer {
	code = stringx.FirstNonEmpty(code)
	return ApplicationOffer{
		ApplicationKey:     code,
		UpstreamServiceKey: code,
		DisplayName:        stringx.FirstNonEmpty(name, code),
	}
}
