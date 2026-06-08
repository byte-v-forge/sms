package httpclient

import (
	"context"
	"net"
	"net/http"
	"net/url"

	xproxy "golang.org/x/net/proxy"
)

func configureSocksProxy(transport *http.Transport, parsed *url.URL) error {
	var auth *xproxy.Auth
	if parsed.User != nil {
		password, _ := parsed.User.Password()
		auth = &xproxy.Auth{User: parsed.User.Username(), Password: password}
	}
	dialer, err := xproxy.SOCKS5("tcp", parsed.Host, auth, xproxy.Direct)
	if err != nil {
		return err
	}
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	return nil
}
