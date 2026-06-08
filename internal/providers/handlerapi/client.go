package handlerapi

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

type HTTPDoer = providerhttp.HTTPDoer

type Client struct {
	endpoint   url.URL
	apiKey     string
	httpClient HTTPDoer
	userAgent  string
}

func New(rawEndpoint, apiKey string, httpClient HTTPDoer) (*Client, error) {
	rawEndpoint = strings.TrimSpace(rawEndpoint)
	if rawEndpoint == "" {
		return nil, core.NewError(core.CodeValidationFailed, "handler api endpoint is required", false)
	}
	endpoint, err := url.Parse(rawEndpoint)
	if err != nil {
		return nil, core.NewError(core.CodeValidationFailed, "invalid handler api endpoint", false)
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, core.NewError(core.CodeValidationFailed, "handler api key is required", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		endpoint:   *endpoint,
		apiKey:     apiKey,
		httpClient: httpClient,
		userAgent:  "sms/1.0",
	}, nil
}
