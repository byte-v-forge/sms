package httpclient

import "strings"

type ConfigError struct {
	Field string
	Msg   string
}

func (e *ConfigError) Error() string {
	if e == nil {
		return ""
	}
	if e.Field == "" {
		return e.Msg
	}
	return e.Field + ": " + e.Msg
}

func IsRetryableTransportError(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	for _, hint := range []string{"tls", "connection reset", "connection aborted", "connection refused", "timed out", "timeout", "temporarily unavailable", "network is unreachable", "proxy", "proxyconnect", "eof"} {
		if strings.Contains(text, hint) {
			return true
		}
	}
	return false
}
