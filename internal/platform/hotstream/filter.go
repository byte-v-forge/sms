package hotstream

import observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"

type Filter struct {
	EventTypes     []string
	SourceServices []string
	ResourceTypes  []string
	ResourceIDs    []string
	Scopes         []string
	Attributes     map[string]string
}

func (f Filter) Match(event *observabilityv1.HotStreamEvent) bool {
	if event == nil {
		return false
	}
	metadata := event.GetMetadata()
	return matchAny(f.EventTypes, metadata.GetType()) &&
		matchAny(f.SourceServices, metadata.GetSource()) &&
		matchAny(f.ResourceTypes, event.GetResourceType()) &&
		matchAny(f.ResourceIDs, event.GetResourceId()) &&
		matchAny(f.Scopes, event.GetScope()) &&
		matchAttributes(f.Attributes, event.GetAttributes())
}
