package herosms

import "context"

func (c *Client) catalogServiceNames(ctx context.Context) map[string]string {
	services, err := c.ListServices(ctx)
	if err != nil {
		return nil
	}
	return heroSMSServiceNameIndex(services)
}
