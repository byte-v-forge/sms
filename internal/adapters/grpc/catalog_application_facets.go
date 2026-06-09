package grpcadapter

import (
	"sort"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func applicationFacets(offers []core.RouteOffer) []*smsv1.SmsApplicationInfo {
	items := map[string]*smsv1.SmsApplicationInfo{}
	for _, offer := range offers {
		if offer.ApplicationKey == "" {
			continue
		}
		items[offer.ApplicationKey] = &smsv1.SmsApplicationInfo{ApplicationKey: offer.ApplicationKey, DisplayName: offer.ApplicationName}
	}
	out := applicationFacetValues(items)
	sort.Slice(out, func(i, j int) bool { return out[i].GetApplicationKey() < out[j].GetApplicationKey() })
	return out
}

func applicationFacetValues(items map[string]*smsv1.SmsApplicationInfo) []*smsv1.SmsApplicationInfo {
	out := make([]*smsv1.SmsApplicationInfo, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	return out
}
