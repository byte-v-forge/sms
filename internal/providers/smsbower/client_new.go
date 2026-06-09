package smsbower

import "github.com/byte-v-forge/sms/internal/providers/handlerapi"

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	api, err := handlerapi.New(config.endpoint(), config.APIKey, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{api: api, policy: defaultProviderPolicy()}, nil
}
