package hotstream

import (
	"strings"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

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

func matchAny(allowed []string, value string) bool {
	if len(allowed) == 0 {
		return true
	}
	value = strings.TrimSpace(value)
	for _, item := range allowed {
		if strings.TrimSpace(item) == value {
			return true
		}
	}
	return false
}

func matchAttributes(expected map[string]string, actual map[string]string) bool {
	if len(expected) == 0 {
		return true
	}
	for key, value := range expected {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if actual[strings.TrimSpace(key)] != strings.TrimSpace(value) {
			return false
		}
	}
	return true
}
