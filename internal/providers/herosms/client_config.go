package herosms

import "strings"

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

func (c Config) withDefaults() Config {
	if c.Endpoint == "" {
		c.Endpoint = DefaultEndpoint
	}
	c.OpenAPIEndpoint = strings.TrimRight(strings.TrimSpace(c.OpenAPIEndpoint), "/")
	if c.OpenAPIEndpoint == "" {
		c.OpenAPIEndpoint = DefaultOpenAPIEndpoint
	}
	return c
}
