package handlerapi

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetNumberV2(ctx context.Context, request core.ProviderAcquireRequest, config GetNumberV2Config) (string, error) {
	params, err := getNumberV2Params(request, config)
	if err != nil {
		return "", err
	}
	return c.Do(ctx, "getNumberV2", params)
}
