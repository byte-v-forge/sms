package herosms

import (
	"net/url"
	"strings"
)

func (c *Client) openAPIEndpointWithPath(path string) url.URL {
	endpoint := c.openAPIEndpoint
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + path
	return endpoint
}
