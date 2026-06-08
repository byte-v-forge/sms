package httpclient

import (
	"net/http"
	"net/url"
	"strings"
)

func configureProxyTransport(transport *http.Transport, parsed *url.URL, allowed map[string]bool) (*http.Transport, error) {
	scheme := strings.ToLower(parsed.Scheme)
	switch scheme {
	case "http", "https":
		if !allowed[scheme] {
			return nil, unsupportedProxyScheme(scheme, allowed)
		}
		transport.Proxy = http.ProxyURL(parsed)
	case "socks5", "socks5h":
		if !allowed[scheme] {
			return nil, unsupportedProxyScheme(scheme, allowed)
		}
		if err := configureSocksProxy(transport, parsed); err != nil {
			return nil, err
		}
	default:
		return nil, unsupportedProxyScheme(scheme, allowed)
	}
	return transport, nil
}
