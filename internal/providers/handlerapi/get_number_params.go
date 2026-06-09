package handlerapi

import (
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func getNumberV2Params(request core.ProviderAcquireRequest, config GetNumberV2Config) (url.Values, error) {
	config = normalizeGetNumberV2Config(config)
	service := strings.TrimSpace(request.Route.UpstreamServiceKey)
	if service == "" {
		return nil, core.NewError(core.CodeValidationFailed, config.ProviderName+" service is required", false)
	}
	country := strings.TrimSpace(request.Route.ProviderCountryID)
	if country == "" {
		return nil, core.NewError(core.CodeValidationFailed, config.ProviderName+" "+config.CountryLabel+" is required", false)
	}
	providerID := strings.TrimSpace(request.Route.UpstreamProviderID)
	if config.ProviderIDRequired && providerID == "" {
		return nil, core.NewError(core.CodeValidationFailed, config.ProviderName+" "+config.ProviderIDLabel+" is required", false)
	}
	return newGetNumberV2Params(request, config, service, country, providerID), nil
}

func newGetNumberV2Params(request core.ProviderAcquireRequest, config GetNumberV2Config, service string, country string, providerID string) url.Values {
	params := url.Values{}
	params.Set("service", service)
	params.Set("country", country)
	if providerID != "" && config.ProviderIDParam != "" {
		params.Set(config.ProviderIDParam, providerID)
	}
	if maxPrice := strings.TrimSpace(request.Route.MaxPrice.AmountDecimal); maxPrice != "" && config.MaxPriceParam != "" {
		params.Set(config.MaxPriceParam, maxPrice)
	}
	return params
}
