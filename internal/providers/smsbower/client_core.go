package smsbower

import (
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

const (
	DefaultEndpoint = "https://smsbower.page/stubs/handler_api.php"
	ProviderKey     = "smsbower"
)

type Client struct {
	api    *handlerapi.Client
	policy core.ProviderPolicy
}

func (c *Client) Key() string {
	return ProviderKey
}

func (c *Client) Policy() core.ProviderPolicy {
	return c.policy
}
