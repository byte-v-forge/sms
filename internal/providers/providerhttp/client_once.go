package providerhttp

import (
	"context"

	"github.com/byte-v-forge/sms/internal/platform/httpx"
)

func doOnce(ctx context.Context, doer HTTPDoer, factory RequestFactory, maxBodyBytes int64) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	req, err := factory(ctx)
	if err != nil {
		return Response{}, err
	}
	resp, err := doer.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := httpx.ReadLimited(resp.Body, maxBodyBytes)
	if err != nil {
		return Response{StatusCode: resp.StatusCode, Header: resp.Header}, err
	}
	return Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: body}, nil
}
