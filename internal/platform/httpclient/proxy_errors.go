package httpclient

import (
	"fmt"
	"sort"
	"strings"
)

func unsupportedProxyScheme(scheme string, allowed map[string]bool) error {
	items := make([]string, 0, len(allowed))
	for item := range allowed {
		items = append(items, item)
	}
	sort.Strings(items)
	return &ConfigError{Field: "proxy_url", Msg: fmt.Sprintf("unsupported proxy scheme %q; supported schemes: %s", scheme, strings.Join(items, ", "))}
}
