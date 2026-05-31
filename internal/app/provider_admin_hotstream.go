package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/hotstream"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminService) publishProviderConfig(ctx context.Context, eventType string, config *smsinternalv1.SmsProviderConfig) {
	if config == nil {
		return
	}
	s.publishResource(ctx, eventType, SMSProviderConfigResource, config.GetProviderKey(), map[string]string{
		"provider_key": config.GetProviderKey(),
		"enabled":      fmt.Sprintf("%t", config.GetEnabled()),
	})
}

func (s *ProviderAdminService) publishResource(ctx context.Context, eventType string, resourceType string, resourceID string, attrs map[string]string) {
	if s == nil || s.hot == nil || strings.TrimSpace(resourceID) == "" {
		return
	}
	now := time.Now()
	event := hotstream.NewEvent(hotstream.EventConfig{
		EventID:       eventbus.StableEventID("sms-hot-", eventType, resourceID, fmt.Sprintf("%d", now.UnixNano())),
		EventType:     eventType,
		SourceService: SMSHotStreamSource,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		OccurredAt:    now,
		CorrelationID: resourceID,
		Attributes:    attrs,
	})
	if err := s.hot.Publish(context.WithoutCancel(ctx), event); err != nil {
		log.Printf("publish sms config hotstream failed type=%s resource=%s: %v", eventType, resourceID, err)
	}
}
