package herosms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/byte-v-forge/common-lib/httpx"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func (c *Client) getOpenAPIJSON(ctx context.Context, path string, params url.Values, out any) error {
	endpoint, err := url.Parse(c.openAPIEndpoint + path)
	if err != nil {
		return core.NewError(core.CodeValidationFailed, "invalid hero sms openapi endpoint", false)
	}
	if len(params) > 0 {
		endpoint.RawQuery = params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return core.NewError(core.CodeInternal, err.Error(), false)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	defer resp.Body.Close()
	body, err := httpx.ReadLimited(resp.Body, 8<<20)
	if err != nil {
		return core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(body))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return mapHeroSMSOpenAPIError(resp.StatusCode, text)
	}
	if err := json.Unmarshal(body, out); err != nil {
		return mapHeroSMSOpenAPIError(resp.StatusCode, text)
	}
	return nil
}

func mapHeroSMSOpenAPIError(statusCode int, text string) error {
	var payload struct {
		Title   string `json:"title"`
		Details string `json:"details"`
	}
	if err := json.Unmarshal([]byte(text), &payload); err == nil && strings.TrimSpace(payload.Title) != "" {
		return handlerapi.MapTextError(strings.TrimSpace(payload.Title))
	}
	if text != "" {
		return handlerapi.MapTextError(text)
	}
	if statusCode >= 500 {
		return core.NewError(core.CodeSupplyUnavailable, fmt.Sprintf("hero sms openapi http status %d", statusCode), true)
	}
	return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("hero sms openapi http status %d", statusCode), false)
}
