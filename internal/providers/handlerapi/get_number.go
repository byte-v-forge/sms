package handlerapi

import (
	"context"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetNumberV2(ctx context.Context, request core.ProviderAcquireRequest, config GetNumberV2Config) (string, error) {
	params, err := getNumberV2Params(request, config)
	if err != nil {
		return "", err
	}
	return c.Do(ctx, "getNumberV2", params)
}

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
	params := url.Values{}
	params.Set("service", service)
	params.Set("country", country)
	if providerID != "" && config.ProviderIDParam != "" {
		params.Set(config.ProviderIDParam, providerID)
	}
	if maxPrice := strings.TrimSpace(request.Route.MaxPrice.AmountDecimal); maxPrice != "" && config.MaxPriceParam != "" {
		params.Set(config.MaxPriceParam, maxPrice)
	}
	return params, nil
}

func normalizeGetNumberV2Config(config GetNumberV2Config) GetNumberV2Config {
	config.ProviderName = strings.TrimSpace(config.ProviderName)
	if config.ProviderName == "" {
		config.ProviderName = "sms provider"
	}
	config.CountryLabel = strings.TrimSpace(config.CountryLabel)
	if config.CountryLabel == "" {
		config.CountryLabel = "country"
	}
	config.ProviderIDLabel = strings.TrimSpace(config.ProviderIDLabel)
	if config.ProviderIDLabel == "" {
		config.ProviderIDLabel = "upstream provider id"
	}
	return config
}
