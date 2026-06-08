package fivesim

import "strings"

const (
	DefaultEndpoint = "https://5sim.net"
	ProviderKey     = "5sim"
)

type Config struct {
	Endpoint     string
	Token        string
	CurrencyCode string
}

func (c Config) withDefaults() Config {
	c.Endpoint = strings.TrimRight(strings.TrimSpace(c.Endpoint), "/")
	if c.Endpoint == "" {
		c.Endpoint = DefaultEndpoint
	}
	return c
}
