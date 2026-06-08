package fivesim

import (
	"net/url"
	"strings"
)

func (c *Client) endpointWithPath(path string) url.URL {
	endpoint := c.endpoint
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + path
	return endpoint
}
