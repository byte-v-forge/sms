package fivesim

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

type Client struct {
	endpoint     url.URL
	token        string
	currencyCode string
	httpClient   handlerapi.HTTPDoer
	userAgent    string
	policy       core.ProviderPolicy
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	config = config.withDefaults()
	endpoint, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, core.NewError(core.CodeValidationFailed, "invalid 5sim endpoint", false)
	}
	if strings.TrimSpace(config.Token) == "" {
		return nil, core.NewError(core.CodeValidationFailed, "5sim token is required", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		endpoint:     *endpoint,
		token:        strings.TrimSpace(config.Token),
		currencyCode: strings.TrimSpace(config.CurrencyCode),
		httpClient:   httpClient,
		userAgent:    "sms/1.0",
		policy:       defaultProviderPolicy(),
	}, nil
}
