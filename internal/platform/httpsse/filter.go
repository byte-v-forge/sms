package httpsse

import (
	"net/http"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func FilterFromRequest(r *http.Request, base hotstream.Filter) hotstream.Filter {
	if r == nil {
		return base
	}
	q := r.URL.Query()
	base.EventTypes = append(base.EventTypes, splitValues(q["event_type"])...)
	base.ResourceTypes = append(base.ResourceTypes, splitValues(q["resource_type"])...)
	base.ResourceIDs = append(base.ResourceIDs, splitValues(q["resource_id"])...)
	base.Scopes = append(base.Scopes, splitValues(q["scope"])...)
	attrs := map[string]string{}
	for key, values := range q {
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
	if len(attrs) > 0 {
		if base.Attributes == nil {
			base.Attributes = attrs
		} else {
			for key, value := range attrs {
				base.Attributes[key] = value
			}
		}
	}
	return base
}

func splitValues(values []string) []string {
	out := []string{}
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}
