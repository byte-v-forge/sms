package app

import (
	"context"
	"log"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) publishOrder(ctx context.Context, order core.Order) {
	if s == nil || s.hot == nil {
		return
	}
	if err := s.hot.Publish(context.WithoutCancel(ctx), hotStreamOrderEvent(order)); err != nil {
		log.Printf("publish sms order hotstream failed order=%s: %v", order.ID, err)
	}
}
