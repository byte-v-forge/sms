package httpclient

import "strings"

var HTTPProxySchemes = []string{"http", "https"}
var CommonProxySchemes = []string{"http", "https", "socks5", "socks5h"}

func normalizedSchemes(values []string) map[string]bool {
	if len(values) == 0 {
		values = CommonProxySchemes
	}
	out := make(map[string]bool, len(values))
	for _, value := range values {
		if value = strings.ToLower(strings.TrimSpace(value)); value != "" {
			out[value] = true
		}
	}
	return out
}
