package herosms

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

type Client struct {
	api             *handlerapi.Client
	apiKey          string
	openAPIEndpoint url.URL
	httpClient      handlerapi.HTTPDoer
	userAgent       string
	policy          core.ProviderPolicy
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	config = config.withDefaults()
	openAPIEndpoint, err := url.Parse(config.OpenAPIEndpoint)
	if err != nil {
		return nil, core.NewError(core.CodeValidationFailed, "invalid hero sms openapi endpoint", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	api, err := handlerapi.New(config.Endpoint, config.APIKey, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{
		api:             api,
		apiKey:          strings.TrimSpace(config.APIKey),
		openAPIEndpoint: *openAPIEndpoint,
		httpClient:      httpClient,
		userAgent:       "sms/1.0",
		policy:          defaultProviderPolicy(),
	}, nil
}
