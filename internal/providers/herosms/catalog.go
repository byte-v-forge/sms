package herosms

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

type ApplicationOffer struct {
	ApplicationKey     string
	UpstreamServiceKey string
	DisplayName        string
}

type Country struct {
	CountryID string
	Name      string
	NameCN    string
	Visible   bool
	Retry     bool
}

type PriceOffer struct {
	CountryID          string
	UpstreamServiceKey string
	Price              core.Money
	AvailableCount     int
	PhysicalCount      int
}

func (c *Client) ListApplications(ctx context.Context) ([]ApplicationOffer, error) {
	result, err := c.api.Do(ctx, "getServicesList", nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Status   string `json:"status"`
		Services []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"services"`
	}
	if err := decodeJSONObject(result, &payload); err != nil {
		return nil, err
	}
	if payload.Status != "" && payload.Status != "success" {
		return nil, handlerapi.MapTextError(payload.Status)
	}
	offers := make([]ApplicationOffer, 0, len(payload.Services))
	for _, service := range payload.Services {
		offers = append(offers, ApplicationOffer{
			ApplicationKey:     service.Code,
			UpstreamServiceKey: service.Code,
			DisplayName:        firstNonEmpty(service.Name, service.Code),
		})
	}
	return offers, nil
}

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	result, err := c.api.Do(ctx, "getCountries", nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID      int    `json:"id"`
		Rus     string `json:"rus"`
		Eng     string `json:"eng"`
		Chn     string `json:"chn"`
		Visible int    `json:"visible"`
		Retry   int    `json:"retry"`
	}
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	countries := make([]Country, 0, len(raw))
	for _, item := range raw {
		countries = append(countries, Country{
			CountryID: strconv.Itoa(item.ID),
			Name:      firstNonEmpty(item.Eng, item.Chn, item.Rus, strconv.Itoa(item.ID)),
			NameCN:    firstNonEmpty(item.Chn, item.Eng, item.Rus, strconv.Itoa(item.ID)),
			Visible:   item.Visible == 1,
			Retry:     item.Retry == 1,
		})
	}
	return countries, nil
}

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	if serviceKey != "" {
		params.Set("service", serviceKey)
	}
	if countryID != "" {
		params.Set("country", countryID)
	}
	result, err := c.api.Do(ctx, "getPrices", params)
	if err != nil {
		return nil, err
	}
	var raw map[string]map[string]struct {
		Cost          json.RawMessage `json:"cost"`
		Count         int             `json:"count"`
		PhysicalCount int             `json:"physicalCount"`
	}
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	var offers []PriceOffer
	for cID, byService := range raw {
		for service, offer := range byService {
			offers = append(offers, PriceOffer{
				CountryID:          cID,
				UpstreamServiceKey: service,
				Price:              core.Money{AmountDecimal: rawJSONScalar(offer.Cost)},
				AvailableCount:     offer.Count,
				PhysicalCount:      offer.PhysicalCount,
			})
		}
	}
	return offers, nil
}

func decodeJSONObject(result string, out any) error {
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return core.NewError(core.CodeUpstreamRejected, "bad json response: "+err.Error(), false)
	}
	return nil
}

func rawJSONScalar(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return number.String()
	}
	var floatValue float64
	if err := json.Unmarshal(raw, &floatValue); err == nil {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	return string(raw)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
