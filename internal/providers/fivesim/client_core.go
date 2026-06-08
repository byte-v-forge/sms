package fivesim

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

const (
	DefaultEndpoint = "https://5sim.net"
	ProviderKey     = "5sim"
)

type Config struct {
	Endpoint     string
	Token        string
	CurrencyCode string
}

type Client struct {
	endpoint     url.URL
	token        string
	currencyCode string
	httpClient   handlerapi.HTTPDoer
	userAgent    string
	policy       core.ProviderPolicy
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	rawEndpoint := strings.TrimRight(strings.TrimSpace(config.Endpoint), "/")
	if rawEndpoint == "" {
		rawEndpoint = DefaultEndpoint
	}
	endpoint, err := url.Parse(rawEndpoint)
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
		policy: core.ProviderPolicy{
			OrderTTL:     20 * time.Minute,
			PollInterval: 5 * time.Second,
		},
	}, nil
}

func (c *Client) Key() string {
	return ProviderKey
}

func (c *Client) Policy() core.ProviderPolicy {
	return c.policy
}
