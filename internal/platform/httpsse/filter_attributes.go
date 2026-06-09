package httpsse

import (
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func filterAttributesFromQuery(query url.Values) map[string]string {
	attrs := map[string]string{}
	for key, values := range query {
		if !strings.HasPrefix(key, "attr.") {
			continue
		}
		name := strings.TrimSpace(strings.TrimPrefix(key, "attr."))
		if name == "" || len(values) == 0 {
			continue
		}
		if value := strings.TrimSpace(values[len(values)-1]); value != "" {
			attrs[name] = value
		}
	}
	return attrs
}

func mergeFilterAttributes(base hotstream.Filter, attrs map[string]string) hotstream.Filter {
	if len(attrs) == 0 {
		return base
	}
	if base.Attributes == nil {
		base.Attributes = attrs
		return base
	}
	for key, value := range attrs {
		base.Attributes[key] = value
	}
	return base
}
