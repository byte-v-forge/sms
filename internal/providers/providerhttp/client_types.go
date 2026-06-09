package providerhttp

import (
	"context"
	"net/http"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestFactory func(context.Context) (*http.Request, error)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}
