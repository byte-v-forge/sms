package smsbower

type Config struct {
	Endpoint string
	APIKey   string
}

func (c Config) endpoint() string {
	if c.Endpoint == "" {
		return DefaultEndpoint
	}
	return c.Endpoint
}
