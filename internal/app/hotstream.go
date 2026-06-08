package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
)

const (
	SMSHotStreamSource        = "sms-service"
	SMSOrderResource          = "sms.order"
	SMSProviderConfigResource = "sms.provider_config"
	SMSOrderUpdatedEvent      = "sms.order.updated"
	SMSProviderConfigUpdated  = "sms.provider_config.updated"
	SMSProviderConfigDeleted  = "sms.provider_config.deleted"
)

type HotStreamPublisher = hotstream.Publisher

func (s *OrderService) saveOrder(ctx context.Context, order core.Order, records ...eventoutbox.Record) error {
	if err := s.store.Save(ctx, order, records...); err != nil {
		return err
	}
	s.publishOrder(ctx, order)
	return nil
}

func (s *OrderService) updateOrder(ctx context.Context, order core.Order, records ...eventoutbox.Record) error {
	if err := s.store.Update(ctx, order, records...); err != nil {
		return err
	}
	s.publishOrder(ctx, order)
	return nil
}

func (s *OrderService) recordCode(ctx context.Context, order core.Order, code core.SMSCode, records ...eventoutbox.Record) error {
	if err := s.store.RecordCode(ctx, order, code, records...); err != nil {
		return err
	}
	s.publishOrder(ctx, order)
	return nil
}

func (s *OrderService) prepareCodeSecret(ctx context.Context, order core.Order, code core.SMSCode) (core.SMSCode, error) {
	if secretref.Configured(code.SecretRef) {
		return core.SMSCode{ReceivedAt: code.ReceivedAt, SecretRef: code.SecretRef}, nil
	}
	if s == nil || s.codeSecrets == nil {
		return core.SMSCode{}, core.NewError(core.CodeInternal, "sms code secret store is not configured", false)
	}
	return s.codeSecrets.Save(ctx, order, code)
}

func (s *OrderService) publishOrder(ctx context.Context, order core.Order) {
	if s == nil || s.hot == nil {
		return
	}
	updatedAt := order.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}
	event := hotstream.NewEvent(hotstream.EventConfig{
		EventID:       eventbus.StableEventID("sms-order-", order.ID, string(order.Status), fmt.Sprintf("%d", updatedAt.UnixNano())),
		EventType:     SMSOrderUpdatedEvent,
		SourceService: SMSHotStreamSource,
		ResourceType:  SMSOrderResource,
		ResourceID:    order.ID,
		Scope:         string(order.Status),
		OccurredAt:    updatedAt,
		CorrelationID: order.RequestID,
		Attributes: map[string]string{
			"order_id":     order.ID,
			"request_id":   order.RequestID,
			"provider_key": order.ProviderKey,
			"status":       string(order.Status),
		},
	})
	if err := s.hot.Publish(context.WithoutCancel(ctx), event); err != nil {
		log.Printf("publish sms order hotstream failed order=%s: %v", order.ID, err)
	}
}
