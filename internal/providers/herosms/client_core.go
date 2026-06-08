package herosms

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

const (
	DefaultEndpoint        = "https://hero-sms.com/stubs/handler_api.php"
	DefaultOpenAPIEndpoint = "https://hero-sms.com/api/v1"
	ProviderKey            = "herosms"
)

type Config struct {
	Endpoint        string
	OpenAPIEndpoint string
	APIKey          string
}

type Client struct {
	api             *handlerapi.Client
	apiKey          string
	openAPIEndpoint url.URL
	httpClient      handlerapi.HTTPDoer
	userAgent       string
	policy          core.ProviderPolicy
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	rawOpenAPIEndpoint := strings.TrimRight(strings.TrimSpace(config.OpenAPIEndpoint), "/")
	if rawOpenAPIEndpoint == "" {
		rawOpenAPIEndpoint = DefaultOpenAPIEndpoint
	}
	openAPIEndpoint, err := url.Parse(rawOpenAPIEndpoint)
	if err != nil {
		return nil, core.NewError(core.CodeValidationFailed, "invalid hero sms openapi endpoint", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	api, err := handlerapi.New(endpoint, config.APIKey, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{
		api:             api,
		apiKey:          strings.TrimSpace(config.APIKey),
		openAPIEndpoint: *openAPIEndpoint,
		httpClient:      httpClient,
		userAgent:       "sms/1.0",
		policy: core.ProviderPolicy{
			OrderTTL:           20 * time.Minute,
			PollInterval:       5 * time.Second,
			CancelAllowedAfter: 2 * time.Minute,
		},
	}, nil
}

func (c *Client) Key() string {
	return ProviderKey
}

func (c *Client) Policy() core.ProviderPolicy {
	return c.policy
}
