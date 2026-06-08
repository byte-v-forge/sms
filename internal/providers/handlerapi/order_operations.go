package handlerapi

import (
	"context"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

type GetNumberV2Config struct {
	ProviderName       string
	CountryLabel       string
	ProviderIDParam    string
	ProviderIDLabel    string
	ProviderIDRequired bool
	MaxPriceParam      string
}

type StatusParser func(string) (core.ProviderCodeResult, error)

type StatusActionMapper func(core.ProviderAction) (status string, expected string, err error)

func (c *Client) GetNumberV2(ctx context.Context, request core.ProviderAcquireRequest, config GetNumberV2Config) (string, error) {
	params, err := getNumberV2Params(request, config)
	if err != nil {
		return "", err
	}
	return c.Do(ctx, "getNumberV2", params)
}

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string, parse StatusParser) (core.ProviderCodeResult, error) {
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	result, err := c.Do(ctx, "getStatus", params)
	if err != nil {
		return core.ProviderCodeResult{}, err
	}
	return parse(result)
}

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction, mapAction StatusActionMapper) error {
	status, expected, err := mapAction(action)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	params.Set("status", status)
	result, err := c.Do(ctx, "setStatus", params)
	if err != nil {
		return err
	}
	if result != expected {
		return MapTextError(result)
	}
	return nil
}

func (c *Client) SetActivationStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction, providerName string) error {
	return c.SetStatus(ctx, upstreamOrderID, action, func(action core.ProviderAction) (string, string, error) {
		return ActivationStatusForAction(providerName, action)
	})
}

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	result, err := c.Do(ctx, "getBalance", nil)
	if err != nil {
		return core.Money{}, err
	}
	const prefix = "ACCESS_BALANCE:"
	if !strings.HasPrefix(result, prefix) {
		return core.Money{}, MapTextError(result)
	}
	return core.Money{AmountDecimal: strings.TrimPrefix(result, prefix)}, nil
}

func ActivationStatusForAction(providerName string, action core.ProviderAction) (status string, expected string, err error) {
	switch action {
	case core.ActionMarkMessageSent:
		return "1", "ACCESS_READY", nil
	case core.ActionRequestAdditional:
		return "3", "ACCESS_RETRY_GET", nil
	case core.ActionCompleteOrder:
		return "6", "ACCESS_ACTIVATION", nil
	case core.ActionCancelOrder:
		return "8", "ACCESS_CANCEL", nil
	default:
		providerName = strings.TrimSpace(providerName)
		if providerName == "" {
			providerName = "sms provider"
		}
		return "", "", core.NewError(core.CodeUnsupportedOperation, "unsupported "+providerName+" status action", false)
	}
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
