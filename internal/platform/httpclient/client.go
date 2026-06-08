package httpclient

import (
	"net/http"
	"time"
)

func New(timeout time.Duration, proxyRawURL string) (*http.Client, error) {
	return NewWithSchemes(timeout, proxyRawURL, CommonProxySchemes...)
}

func NewWithSchemes(timeout time.Duration, proxyRawURL string, schemes ...string) (*http.Client, error) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	transport, err := Transport(proxyRawURL, schemes...)
	if err != nil {
		return nil, err
	}
	return &http.Client{Timeout: timeout, Transport: transport}, nil
}
