package httpsse

import (
	"net/http"

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
	return mergeFilterAttributes(base, filterAttributesFromQuery(q))
}
