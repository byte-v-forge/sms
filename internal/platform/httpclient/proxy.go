package httpclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Transport(proxyRawURL string, schemes ...string) (*http.Transport, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	proxyRawURL = strings.TrimSpace(proxyRawURL)
	if proxyRawURL == "" {
		return transport, nil
	}
	parsed, err := url.Parse(proxyRawURL)
	if err != nil {
		return nil, fmt.Errorf("parse proxy_url: %w", err)
	}
	return configureProxyTransport(transport, parsed, normalizedSchemes(schemes))
}
