package smsbower

import (
	"context"
	"encoding/json"
	"github.com/byte-v-forge/common-lib/stringx"
	"net/url"
	"strings"
	"time"

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
}

type PriceOffer struct {
	CountryID          string
	UpstreamServiceKey string
	ProviderID         string
	Price              core.Money
	AvailableCount     int
}

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

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	result, err := c.api.Do(ctx, "getCountries", nil)
	if err != nil {
		return nil, err
	}
	var raw map[string]struct {
		ID  json.RawMessage `json:"id"`
		Rus string          `json:"rus"`
		Eng string          `json:"eng"`
		Chn string          `json:"chn"`
	}
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	countries := make([]Country, 0, len(raw))
	for key, item := range raw {
		id := rawJSONScalar(item.ID)
		if id == "" || id == "0" {
			id = key
		}
		name := stringx.FirstNonEmpty(item.Eng, item.Chn, item.Rus, key)
		countries = append(countries, Country{CountryID: id, Name: name})
	}
	return countries, nil
}

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	params.Set("service", serviceKey)
	params.Set("country", countryID)
	result, err := c.api.Do(ctx, "getPricesV3", params)
	if err != nil {
		return nil, err
	}
	var raw map[string]map[string]map[string]struct {
		Count      int             `json:"count"`
		Price      json.RawMessage `json:"price"`
		ProviderID json.RawMessage `json:"provider_id"`
	}
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	var offers []PriceOffer
	for cID, byService := range raw {
		for svc, byProvider := range byService {
			for providerID, offer := range byProvider {
				offers = append(offers, PriceOffer{
					CountryID:          cID,
					UpstreamServiceKey: svc,
					ProviderID:         stringx.FirstNonEmpty(rawJSONScalar(offer.ProviderID), providerID),
					Price:              core.Money{AmountDecimal: rawJSONScalar(offer.Price)},
					AvailableCount:     offer.Count,
				})
			}
		}
	}
	return offers, nil
}

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	applications, err := c.ListApplications(ctx)
	if err != nil {
		return nil, err
	}
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	service := matchService(query.ApplicationKey, applications)
	if strings.TrimSpace(query.ApplicationKey) != "" && service == "" {
		return nil, nil
	}
	countryID := smsbowerCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return nil, nil
	}
	priceOffers, err := c.ListPriceOffers(ctx, service, countryID)
	if err != nil {
		return nil, err
	}
	appNames := map[string]string{}
	for _, app := range applications {
		appNames[app.UpstreamServiceKey] = stringx.FirstNonEmpty(app.DisplayName, app.UpstreamServiceKey)
	}
	countryByID := map[string]Country{}
	for _, country := range countries {
		countryByID[country.CountryID] = country
	}
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		country := countryByID[offer.CountryID]
		iso2, callingCode := smsbowerCountryCodes(country.Name)
		route := core.Route{
			ProviderKey:        ProviderKey,
			ApplicationKey:     offer.UpstreamServiceKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			CountryISO2:        iso2,
			CountryCallingCode: callingCode,
			ProviderCountryID:  offer.CountryID,
			UpstreamProviderID: offer.ProviderID,
		}
		out = append(out, core.RouteOffer{
			ProviderKey:          ProviderKey,
			UpstreamProviderID:   offer.ProviderID,
			UpstreamProviderName: offer.ProviderID,
			ApplicationKey:       offer.UpstreamServiceKey,
			ApplicationName:      stringx.FirstNonEmpty(appNames[offer.UpstreamServiceKey], offer.UpstreamServiceKey),
			CountryISO2:          iso2,
			CountryName:          stringx.FirstNonEmpty(country.Name, offer.CountryID),
			CountryCallingCode:   callingCode,
			Price:                offer.Price,
			AvailableCount:       offer.AvailableCount,
			ObservedAt:           now,
			Route:                route,
		})
	}
	return out, nil
}

func smsbowerCountryIDForQuery(query core.RouteOfferQuery, countries []Country) string {
	for _, country := range countries {
		iso2, callingCode := smsbowerCountryCodes(country.Name)
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, iso2) {
			continue
		}
		if query.CountryCallingCode != "" && query.CountryCallingCode != strings.TrimPrefix(callingCode, "+") {
			continue
		}
		return country.CountryID
	}
	return ""
}

func smsbowerCountryCodes(name string) (string, string) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	switch {
	case strings.Contains(normalized, "indonesia"):
		return "ID", "62"
	case strings.Contains(normalized, "thai"):
		return "TH", "66"
	default:
		return "", ""
	}
}
