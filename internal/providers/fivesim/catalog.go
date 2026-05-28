package fivesim

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

type Country struct {
	CountryID          string
	Name               string
	CountryISO2        string
	CountryCallingCode string
}

type PriceOffer struct {
	CountryID          string
	UpstreamServiceKey string
	Operator           string
	Price              core.Money
	AvailableCount     int
	SuccessRate        float64
}

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var raw map[string]struct {
		ISO    map[string]int `json:"iso"`
		Prefix map[string]int `json:"prefix"`
		Name   string         `json:"text_en"`
	}
	if err := c.getJSON(ctx, "/v1/guest/countries", nil, false, &raw); err != nil {
		return nil, err
	}
	countries := make([]Country, 0, len(raw))
	for countryID, item := range raw {
		countries = append(countries, Country{
			CountryID:          countryID,
			Name:               item.Name,
			CountryISO2:        firstMapKey(item.ISO),
			CountryCallingCode: trimPlus(firstMapKey(item.Prefix)),
		})
	}
	return countries, nil
}

func (c *Client) ListPriceOffers(ctx context.Context, product, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	if product != "" {
		params.Set("product", product)
	}
	if countryID != "" {
		params.Set("country", countryID)
	}
	var raw map[string]map[string]map[string]struct {
		Cost  json.RawMessage `json:"cost"`
		Count int             `json:"count"`
		Rate  float64         `json:"rate"`
	}
	if err := c.getJSON(ctx, "/v1/guest/prices", params, false, &raw); err != nil {
		return nil, err
	}
	var offers []PriceOffer
	for country, byProduct := range raw {
		for productKey, byOperator := range byProduct {
			for operator, offer := range byOperator {
				offers = append(offers, PriceOffer{
					CountryID:          country,
					UpstreamServiceKey: productKey,
					Operator:           operator,
					Price:              core.Money{CurrencyCode: c.currencyCode, AmountDecimal: rawJSONScalar(offer.Cost)},
					AvailableCount:     offer.Count,
					SuccessRate:        offer.Rate,
				})
			}
		}
	}
	return offers, nil
}

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	countryID := fiveSimCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return nil, nil
	}
	priceOffers, err := c.ListPriceOffers(ctx, query.ApplicationKey, countryID)
	if err != nil {
		return nil, err
	}
	countryByID := map[string]Country{}
	for _, country := range countries {
		countryByID[country.CountryID] = country
	}
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		country := countryByID[offer.CountryID]
		route := core.Route{
			ProviderKey:        ProviderKey,
			ApplicationKey:     offer.UpstreamServiceKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			CountryISO2:        country.CountryISO2,
			CountryCallingCode: country.CountryCallingCode,
			ProviderCountryID:  offer.CountryID,
			UpstreamProviderID: offer.Operator,
		}
		out = append(out, core.RouteOffer{
			ProviderKey:          ProviderKey,
			UpstreamProviderID:   offer.Operator,
			UpstreamProviderName: offer.Operator,
			ApplicationKey:       offer.UpstreamServiceKey,
			ApplicationName:      offer.UpstreamServiceKey,
			CountryISO2:          country.CountryISO2,
			CountryName:          country.Name,
			CountryCallingCode:   country.CountryCallingCode,
			Price:                offer.Price,
			AvailableCount:       offer.AvailableCount,
			ObservedAt:           now,
			Route:                route,
		})
	}
	return out, nil
}

func fiveSimCountryIDForQuery(query core.RouteOfferQuery, countries []Country) string {
	for _, country := range countries {
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, country.CountryISO2) {
			continue
		}
		if query.CountryCallingCode != "" && strings.TrimPrefix(country.CountryCallingCode, "+") != strings.TrimPrefix(query.CountryCallingCode, "+") {
			continue
		}
		return country.CountryID
	}
	return ""
}

func firstMapKey(values map[string]int) string {
	for key := range values {
		return key
	}
	return ""
}

func trimPlus(value string) string {
	if len(value) > 0 && value[0] == '+' {
		return value[1:]
	}
	return value
}
