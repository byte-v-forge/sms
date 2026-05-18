package catalog

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/fivesim"
	"github.com/byte-v-forge/sms/internal/providers/herosms"
	"github.com/byte-v-forge/sms/internal/providers/smsbower"
)

type ProviderConfig struct {
	ProviderKey  string
	APIKey       string
	Endpoint     string
	CurrencyCode string
}

type ListOptions struct {
	ServiceKey string
	CountryID  string
	Operator   string
	MinCount   int
	Limit      int
}

type Catalog struct {
	ProviderKey string        `json:"provider_key"`
	Balance     core.Money    `json:"balance"`
	Countries   []Country     `json:"countries"`
	Services    []Application `json:"services"`
	Options     []Option      `json:"options"`
}

type Country struct {
	ProviderCountryID  string `json:"provider_country_id"`
	Name               string `json:"name"`
	NameCN             string `json:"name_cn,omitempty"`
	CountryISO2        string `json:"country_iso2,omitempty"`
	CountryCallingCode string `json:"country_calling_code,omitempty"`
	RetrySupported     bool   `json:"retry_supported,omitempty"`
}

type Application struct {
	ApplicationKey     string `json:"application_key"`
	UpstreamServiceKey string `json:"upstream_service_key"`
	DisplayName        string `json:"display_name"`
}

type Option struct {
	ProviderKey        string     `json:"provider_key"`
	ApplicationKey     string     `json:"application_key"`
	UpstreamServiceKey string     `json:"upstream_service_key"`
	ProviderCountryID  string     `json:"provider_country_id"`
	CountryName        string     `json:"country_name,omitempty"`
	CountryNameCN      string     `json:"country_name_cn,omitempty"`
	CountryISO2        string     `json:"country_iso2,omitempty"`
	CountryCallingCode string     `json:"country_calling_code,omitempty"`
	ServiceName        string     `json:"service_name,omitempty"`
	UpstreamProviderID string     `json:"upstream_provider_id,omitempty"`
	UpstreamOperator   string     `json:"upstream_operator,omitempty"`
	Price              core.Money `json:"price"`
	AvailableCount     int        `json:"available_count"`
	PhysicalCount      int        `json:"physical_count,omitempty"`
	SuccessRate        float64    `json:"success_rate,omitempty"`
}

func List(ctx context.Context, config ProviderConfig, options ListOptions) (Catalog, error) {
	switch strings.ToLower(strings.TrimSpace(config.ProviderKey)) {
	case fivesim.ProviderKey:
		return listFiveSim(ctx, config, options)
	case herosms.ProviderKey:
		return listHeroSMS(ctx, config, options)
	case smsbower.ProviderKey:
		return listSMSBower(ctx, config, options)
	default:
		return Catalog{}, core.NewError(core.CodeValidationFailed, "unsupported provider key", false)
	}
}

func listFiveSim(ctx context.Context, config ProviderConfig, options ListOptions) (Catalog, error) {
	client, err := fivesim.New(fivesim.Config{
		Endpoint:        config.Endpoint,
		Token:           config.APIKey,
		CurrencyCode:    firstNonEmpty(config.CurrencyCode, "RUB"),
		DefaultOperator: firstNonEmpty(options.Operator, "any"),
	}, nil)
	if err != nil {
		return Catalog{}, err
	}
	balance, err := client.GetBalance(ctx)
	if err != nil {
		return Catalog{}, err
	}
	countries, err := client.ListCountries(ctx)
	if err != nil {
		return Catalog{}, err
	}
	priceOffers, err := client.ListPriceOffers(ctx, options.ServiceKey, options.CountryID)
	if err != nil {
		return Catalog{}, err
	}

	catalog := Catalog{ProviderKey: fivesim.ProviderKey, Balance: balance}
	countryByID := make(map[string]Country, len(countries))
	for _, country := range countries {
		item := Country{
			ProviderCountryID:  country.CountryID,
			Name:               firstNonEmpty(country.Name, country.CountryID),
			CountryISO2:        strings.ToUpper(country.CountryISO2),
			CountryCallingCode: country.CountryCallingCode,
		}
		catalog.Countries = append(catalog.Countries, item)
		countryByID[country.CountryID] = item
	}

	serviceSeen := map[string]bool{}
	for _, offer := range priceOffers {
		if offer.AvailableCount < options.MinCount {
			continue
		}
		country := countryByID[offer.CountryID]
		serviceSeen[offer.UpstreamServiceKey] = true
		catalog.Options = append(catalog.Options, Option{
			ProviderKey:        fivesim.ProviderKey,
			ApplicationKey:     offer.UpstreamServiceKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			ProviderCountryID:  offer.CountryID,
			CountryName:        country.Name,
			CountryISO2:        country.CountryISO2,
			CountryCallingCode: country.CountryCallingCode,
			ServiceName:        offer.UpstreamServiceKey,
			UpstreamProviderID: offer.Operator,
			UpstreamOperator:   offer.Operator,
			Price:              offer.Price,
			AvailableCount:     offer.AvailableCount,
			SuccessRate:        offer.SuccessRate,
		})
	}
	for service := range serviceSeen {
		catalog.Services = append(catalog.Services, Application{ApplicationKey: service, UpstreamServiceKey: service, DisplayName: service})
	}
	return finalize(catalog, options), nil
}

func listHeroSMS(ctx context.Context, config ProviderConfig, options ListOptions) (Catalog, error) {
	client, err := herosms.New(herosms.Config{Endpoint: config.Endpoint, APIKey: config.APIKey}, nil)
	if err != nil {
		return Catalog{}, err
	}
	balance, err := client.GetBalance(ctx)
	if err != nil {
		return Catalog{}, err
	}
	countries, err := client.ListCountries(ctx)
	if err != nil {
		return Catalog{}, err
	}
	services, err := client.ListApplications(ctx)
	if err != nil {
		return Catalog{}, err
	}
	priceOffers, err := client.ListPriceOffers(ctx, options.ServiceKey, options.CountryID)
	if err != nil {
		return Catalog{}, err
	}

	catalog := Catalog{ProviderKey: herosms.ProviderKey, Balance: balance}
	countryByID := make(map[string]Country, len(countries))
	for _, country := range countries {
		item := Country{
			ProviderCountryID: country.CountryID,
			Name:              firstNonEmpty(country.Name, country.CountryID),
			NameCN:            country.NameCN,
			RetrySupported:    country.Retry,
		}
		catalog.Countries = append(catalog.Countries, item)
		countryByID[country.CountryID] = item
	}
	serviceByKey := make(map[string]Application, len(services))
	for _, service := range services {
		item := Application{
			ApplicationKey:     service.ApplicationKey,
			UpstreamServiceKey: service.UpstreamServiceKey,
			DisplayName:        service.DisplayName,
		}
		catalog.Services = append(catalog.Services, item)
		serviceByKey[service.UpstreamServiceKey] = item
	}
	for _, offer := range priceOffers {
		if offer.AvailableCount < options.MinCount {
			continue
		}
		country := countryByID[offer.CountryID]
		service := serviceByKey[offer.UpstreamServiceKey]
		catalog.Options = append(catalog.Options, Option{
			ProviderKey:        herosms.ProviderKey,
			ApplicationKey:     firstNonEmpty(service.ApplicationKey, offer.UpstreamServiceKey),
			UpstreamServiceKey: offer.UpstreamServiceKey,
			ProviderCountryID:  offer.CountryID,
			CountryName:        country.Name,
			CountryNameCN:      country.NameCN,
			ServiceName:        firstNonEmpty(service.DisplayName, offer.UpstreamServiceKey),
			Price:              offer.Price,
			AvailableCount:     offer.AvailableCount,
			PhysicalCount:      offer.PhysicalCount,
		})
	}
	return finalize(catalog, options), nil
}

func listSMSBower(ctx context.Context, config ProviderConfig, options ListOptions) (Catalog, error) {
	client, err := smsbower.New(smsbower.Config{Endpoint: config.Endpoint, APIKey: config.APIKey}, nil)
	if err != nil {
		return Catalog{}, err
	}
	balance, err := client.GetBalance(ctx)
	if err != nil {
		return Catalog{}, err
	}
	countries, err := client.ListCountries(ctx)
	if err != nil {
		return Catalog{}, err
	}
	services, err := client.ListApplications(ctx)
	if err != nil {
		return Catalog{}, err
	}
	priceOffers, err := client.ListPriceOffers(ctx, options.ServiceKey, options.CountryID)
	if err != nil {
		return Catalog{}, err
	}

	catalog := Catalog{ProviderKey: smsbower.ProviderKey, Balance: balance}
	countryByID := make(map[string]Country, len(countries))
	for _, country := range countries {
		item := Country{ProviderCountryID: country.CountryID, Name: firstNonEmpty(country.Name, country.CountryID)}
		catalog.Countries = append(catalog.Countries, item)
		countryByID[country.CountryID] = item
	}
	serviceByKey := make(map[string]Application, len(services))
	for _, service := range services {
		item := Application{
			ApplicationKey:     service.ApplicationKey,
			UpstreamServiceKey: service.UpstreamServiceKey,
			DisplayName:        service.DisplayName,
		}
		catalog.Services = append(catalog.Services, item)
		serviceByKey[service.UpstreamServiceKey] = item
	}
	for _, offer := range priceOffers {
		if offer.AvailableCount < options.MinCount {
			continue
		}
		country := countryByID[offer.CountryID]
		service := serviceByKey[offer.UpstreamServiceKey]
		catalog.Options = append(catalog.Options, Option{
			ProviderKey:        smsbower.ProviderKey,
			ApplicationKey:     firstNonEmpty(service.ApplicationKey, offer.UpstreamServiceKey),
			UpstreamServiceKey: offer.UpstreamServiceKey,
			ProviderCountryID:  offer.CountryID,
			CountryName:        country.Name,
			ServiceName:        firstNonEmpty(service.DisplayName, offer.UpstreamServiceKey),
			UpstreamProviderID: offer.ProviderID,
			Price:              offer.Price,
			AvailableCount:     offer.AvailableCount,
		})
	}
	return finalize(catalog, options), nil
}

func finalize(catalog Catalog, options ListOptions) Catalog {
	sort.Slice(catalog.Countries, func(i, j int) bool {
		return catalog.Countries[i].ProviderCountryID < catalog.Countries[j].ProviderCountryID
	})
	sort.Slice(catalog.Services, func(i, j int) bool {
		return catalog.Services[i].DisplayName < catalog.Services[j].DisplayName
	})
	sort.Slice(catalog.Options, func(i, j int) bool {
		leftScore := catalog.Options[i].PhysicalCount
		rightScore := catalog.Options[j].PhysicalCount
		if leftScore == 0 {
			leftScore = catalog.Options[i].AvailableCount
		}
		if rightScore == 0 {
			rightScore = catalog.Options[j].AvailableCount
		}
		if leftScore != rightScore {
			return leftScore > rightScore
		}
		if catalog.Options[i].ProviderCountryID != catalog.Options[j].ProviderCountryID {
			return catalog.Options[i].ProviderCountryID < catalog.Options[j].ProviderCountryID
		}
		return fmt.Sprintf("%s:%s", catalog.Options[i].UpstreamServiceKey, catalog.Options[i].UpstreamProviderID) <
			fmt.Sprintf("%s:%s", catalog.Options[j].UpstreamServiceKey, catalog.Options[j].UpstreamProviderID)
	})
	if options.Limit > 0 && len(catalog.Options) > options.Limit {
		catalog.Options = catalog.Options[:options.Limit]
	}
	return catalog
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
