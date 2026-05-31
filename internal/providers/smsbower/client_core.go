package smsbower

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

const (
	DefaultEndpoint = "https://smsbower.page/stubs/handler_api.php"
	ProviderKey     = "smsbower"
)

type Config struct {
	Endpoint string
	APIKey   string
}

type Client struct {
	api    *handlerapi.Client
	policy core.ProviderPolicy
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	api, err := handlerapi.New(endpoint, config.APIKey, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{
		api: api,
		policy: core.ProviderPolicy{
			OrderTTL:              25 * time.Minute,
			PollInterval:          5 * time.Second,
			EarlyCancelRetryAfter: 2 * time.Minute,
		},
	}, nil
}

func (c *Client) Key() string {
	return ProviderKey
}

func (c *Client) Policy() core.ProviderPolicy {
	return c.policy
}
