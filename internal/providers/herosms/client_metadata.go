package herosms

import "github.com/byte-v-forge/sms/internal/core"

func (c *Client) Key() string {
	return ProviderKey
}

func (c *Client) Policy() core.ProviderPolicy {
	return c.policy
}
