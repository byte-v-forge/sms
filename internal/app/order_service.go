package app

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/eventoutbox"
	"github.com/byte-v-forge/sms/internal/core"
)

type CancelRetryError struct {
	RetryAt time.Time
}

func (e *CancelRetryError) Error() string {
	if e == nil || e.RetryAt.IsZero() {
		return "sms order cancel retry is scheduled"
	}
	return "sms order cancel retry is scheduled for " + e.RetryAt.UTC().Format(time.RFC3339)
}

type OrderService struct {
	store     OrderStore
	hot       HotStreamPublisher
	providers map[string]core.Provider
	clock     core.Clock
	ids       core.IDGenerator
	events    OrderEventSink
}

type OrderStore interface {
	Save(ctx context.Context, order core.Order, events ...eventoutbox.Record) error
	Get(ctx context.Context, orderID string) (core.Order, error)
	Update(ctx context.Context, order core.Order, events ...eventoutbox.Record) error
}

func NewOrderService(
	store OrderStore,
	providers []core.Provider,
	clock core.Clock,
	ids core.IDGenerator,
	events OrderEventSink,
	hot HotStreamPublisher,
) *OrderService {
	index := make(map[string]core.Provider, len(providers))
	for _, provider := range providers {
		index[provider.Key()] = provider
	}
	if clock == nil {
		clock = SystemClock{}
	}
	if ids == nil {
		ids = RandomIDGenerator{}
	}
	if events == nil {
		events = noopOrderEventSink{}
	}
	return &OrderService{
		store:     store,
		providers: index,
		clock:     clock,
		ids:       ids,
		events:    events,
		hot:       hot,
	}
}

func (s *OrderService) AcquireNumber(ctx context.Context, cmd core.AcquireNumberCommand) (core.Order, error) {
	route := cmd.AcquireParams
	if err := validateAcquireRoute(route); err != nil {
		return core.Order{}, err
	}
	if cmd.RequestID == "" {
		cmd.RequestID = s.ids.NewID("req_")
	}
	now := s.clock.Now()
	target := withRouteTargetDefaults(core.Target{}, route)
	expiresAt := orderRequestExpiresAt(now, s.routePolicy(ctx, route), cmd.LeaseDuration)
	order := core.Order{
		ID:          s.ids.NewID("ord_"),
		RequestID:   cmd.RequestID,
		ProviderKey: route.ProviderKey,
		Target:      target,
		Status:      core.StatusAcquireRequested,
		ExpiresAt:   expiresAt,
		UpdatedAt:   now,
	}
	record, err := s.events.OrderAcquireRequested(ctx, order, route, "api_request")
	if err != nil {
		return core.Order{}, err
	}
	if err := s.saveOrder(ctx, order, record); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func validateAcquireRoute(route core.Route) error {
	if strings.TrimSpace(route.ProviderKey) == "" || strings.TrimSpace(route.UpstreamServiceKey) == "" || strings.TrimSpace(route.ProviderCountryID) == "" {
		return core.NewError(core.CodeValidationFailed, "sms acquire params are incomplete", false)
	}
	switch normalizeProviderKey(route.ProviderKey) {
	case "5sim", "smsbower":
		if strings.TrimSpace(route.UpstreamProviderID) == "" {
			return core.NewError(core.CodeValidationFailed, "sms upstream provider id is required", false)
		}
	}
	return nil
}

func (s *OrderService) RunAcquireRequest(ctx context.Context, orderID string, requestID string, route core.Route) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if order.Status != core.StatusAcquireRequested {
		return order, nil
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, err
	}
	if strings.TrimSpace(route.ProviderKey) == "" {
		route = routeFromOrder(order)
	}
	if strings.TrimSpace(route.ProviderKey) == "" {
		return s.recordAcquireFailure(ctx, order, core.NewError(core.CodeRouteNotFound, "sms acquire route not found", false))
	}
	if err := validateAcquireRoute(route); err != nil {
		return s.recordAcquireFailure(ctx, order, asCoreError(err))
	}
	acquired, acquireErr := s.acquireWithRoute(ctx, order, requestID, route)
	if acquireErr == nil {
		return acquired, nil
	}
	return s.recordAcquireFailure(ctx, order, asCoreError(acquireErr))
}

func (s *OrderService) acquireWithRoute(ctx context.Context, order core.Order, requestID string, route core.Route) (core.Order, error) {
	provider, err := s.provider(route.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	target := withRouteTargetDefaults(order.Target, route)
	providerOrder, err := provider.AcquireNumber(ctx, core.ProviderAcquireRequest{
		RequestID:     firstNonEmpty(requestID, order.RequestID),
		Route:         route,
		Target:        target,
		LeaseDuration: remainingLease(s.clock.Now(), order.ExpiresAt),
	})
	if err != nil {
		return core.Order{}, err
	}
	return s.applyProviderOrder(ctx, order, route, provider, providerOrder)
}

func (s *OrderService) applyProviderOrder(ctx context.Context, order core.Order, route core.Route, provider core.Provider, providerOrder core.ProviderOrder) (core.Order, error) {
	now := s.clock.Now()
	acquiredAt := providerOrder.AcquiredAt
	if acquiredAt.IsZero() {
		acquiredAt = now
	}
	policy := providerPolicyForUpstreamOrder(provider, providerOrder.UpstreamOrderID).WithDefaults()
	expiresAt := providerOrder.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = acquiredAt.Add(policy.OrderTTL)
	}
	if !order.ExpiresAt.IsZero() && order.ExpiresAt.Before(expiresAt) {
		expiresAt = order.ExpiresAt
	}
	var cancelAllowedAt time.Time
	if policy.CancelAllowedAfter > 0 {
		cancelAllowedAt = acquiredAt.Add(policy.CancelAllowedAfter)
	}
	previousStatus := order.Status
	order.ProviderKey = provider.Key()
	order.UpstreamOrderID = providerOrder.UpstreamOrderID
	order.Target = withRouteTargetDefaults(order.Target, route)
	order.PhoneNumber = providerOrder.PhoneNumber
	order.Status = core.StatusPendingCode
	order.Price = providerOrder.Price
	order.AcquiredAt = acquiredAt
	order.ExpiresAt = expiresAt
	order.UpdatedAt = now
	order.CancelAllowedAt = cancelAllowedAt
	order.CanRequestAdditionalCode = providerOrder.CanRequestAdditionalCode
	order.LastError = nil
	records, err := s.orderAcquiredRecords(ctx, order, "order_acquired")
	if err != nil {
		return core.Order{}, err
	}
	statusRecords, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	records = append(records, statusRecords...)
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func (s *OrderService) recordAcquireFailure(ctx context.Context, order core.Order, smsErr *core.Error) (core.Order, error) {
	now := s.clock.Now()
	order.LastError = smsErr
	order.UpdatedAt = now
	if smsErr.Retryable {
		if err := s.updateOrder(ctx, order); err != nil {
			return core.Order{}, err
		}
		return order, smsErr
	}
	previousStatus := order.Status
	order.Status = core.StatusFailed
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, smsErr
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func (s *OrderService) CheckCode(ctx context.Context, orderID string) (core.Order, *core.SMSCode, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return core.Order{}, nil, err
	}
	if order.Status == core.StatusCodeReceived {
		return order, nil, nil
	}
	if order.Status == core.StatusAcquireRequested || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return order, nil, nil
	}
	if order.Status.IsFinal() {
		return order, nil, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	previousStatus := order.Status
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, nil, err
	}
	bindOrderProviderConfig(provider, order)
	result, err := provider.GetStatus(ctx, order.UpstreamOrderID)
	if err != nil {
		order.LastError = asCoreError(err)
		order.UpdatedAt = s.clock.Now()
		_ = s.updateOrder(ctx, order)
		return order, nil, err
	}
	order.UpdatedAt = s.clock.Now()
	switch result.Status {
	case core.StatusCodeReceived:
		receivedAt := result.ReceivedAt
		if receivedAt.IsZero() {
			receivedAt = order.UpdatedAt
		}
		code := &core.SMSCode{Value: result.Code, MessageText: result.MessageText, ReceivedAt: receivedAt}
		order.Status = core.StatusCodeReceived
		records, err := s.statusAndCodeRecords(ctx, order, previousStatus, *code)
		if err != nil {
			return core.Order{}, nil, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, nil, err
		}
		return order, code, nil
	case core.StatusPendingCode, core.StatusAdditionalCodeRequested:
		order.Status = result.Status
	case core.StatusCanceled, core.StatusFailed, core.StatusExpired:
		order.Status = result.Status
	}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, nil, err
	}
	return order, nil, nil
}

func (s *OrderService) PollInterval(ctx context.Context, order core.Order) time.Duration {
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return 5 * time.Second
	}
	return providerPolicyForOrder(ctx, provider, order).WithDefaults().PollInterval
}

func (s *OrderService) MarkMessageSent(ctx context.Context, orderID, requestID string) (core.Order, error) {
	return s.applyAction(ctx, orderID, requestID, core.ActionMarkMessageSent, core.StatusMessageSent)
}

func (s *OrderService) RequestAdditionalCode(ctx context.Context, orderID, requestID string) (core.Order, error) {
	return s.applyAction(ctx, orderID, requestID, core.ActionRequestAdditional, core.StatusAdditionalCodeRequested)
}

func (s *OrderService) CompleteOrder(ctx context.Context, orderID, requestID string) (core.Order, error) {
	return s.applyAction(ctx, orderID, requestID, core.ActionCompleteOrder, core.StatusCompleted)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID, requestID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, err
	}
	if order.Status.IsFinal() {
		return order, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	if orderHasCode(order) {
		return order, nil
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		previousStatus := order.Status
		order.Status = core.StatusCanceled
		order.UpdatedAt = s.clock.Now()
		order.LastError = nil
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, err
		}
		return order, nil
	}
	return s.queueCancelRequest(ctx, order, requestID, "api_request")
}

func (s *OrderService) RunCancelRequest(ctx context.Context, orderID, requestID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	return s.cancelLoadedOrder(ctx, order, requestID)
}

func (s *OrderService) cancelLoadedOrder(ctx context.Context, order core.Order, requestID string) (core.Order, error) {
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	bindOrderProviderConfig(provider, order)
	policy := providerPolicyForOrder(ctx, provider, order).WithDefaults()
	now := s.clock.Now()
	if order.Status.IsFinal() {
		return order, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	if orderHasCode(order) {
		return order, nil
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		previousStatus := order.Status
		order.Status = core.StatusCanceled
		order.UpdatedAt = now
		order.LastError = nil
		order.CancelAllowedAt = time.Time{}
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, err
		}
		return order, nil
	}
	if order.IsExpired(now) {
		previousStatus := order.Status
		order.Status = core.StatusExpired
		order.UpdatedAt = now
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return order, err
		}
		return order, core.NewError(core.CodeOrderExpired, "order expired", false)
	}
	if synced, handled, syncErr := s.syncTerminalOrChargedProviderState(ctx, order, provider); syncErr != nil {
		return order, syncErr
	} else if handled {
		return synced, nil
	}
	if !order.CancelAllowedAt.IsZero() && now.Before(order.CancelAllowedAt) {
		return order, &CancelRetryError{RetryAt: order.CancelAllowedAt}
	}
	age := now.Sub(order.AcquiredAt)
	if policy.CancelAllowedAfter > 0 && age < policy.CancelAllowedAfter {
		return s.deferCancelRetry(ctx, order, requestID, order.AcquiredAt.Add(policy.CancelAllowedAfter))
	}
	if policy.CancelAllowedUntil > 0 && age > policy.CancelAllowedUntil {
		return order, core.NewError(core.CodeCancelNotAllowed, "order is too old to cancel", false)
	}
	if err := provider.SetStatus(ctx, order.UpstreamOrderID, core.ActionCancelOrder); err != nil {
		smsErr := asCoreError(err)
		if raced, ok, raceErr := s.resolveCancelRace(ctx, order, provider, smsErr); raceErr != nil {
			return order, raceErr
		} else if ok {
			return raced, nil
		}
		if shouldQueueEarlyCancelRetry(smsErr, policy) {
			return s.deferCancelRetry(ctx, order, requestID, earlyCancelRetryAt(order, policy, now))
		}
		order.LastError = smsErr
		order.UpdatedAt = now
		if smsErr.Retryable {
			order.CancelAllowedAt = earlyCancelRetryAt(order, policy, now)
		}
		_ = s.updateOrder(ctx, order)
		return order, err
	}
	previousStatus := order.Status
	order.Status = core.StatusCanceled
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func (s *OrderService) applyAction(ctx context.Context, orderID, _ string, action core.ProviderAction, next core.OrderStatus) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, err
	}
	if order.Status.IsFinal() {
		return order, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return order, core.NewError(core.CodeOrderNotFound, "order has no upstream provider lease", true)
	}
	previousStatus := order.Status
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	bindOrderProviderConfig(provider, order)
	if err := provider.SetStatus(ctx, order.UpstreamOrderID, action); err != nil {
		order.LastError = asCoreError(err)
		order.UpdatedAt = s.clock.Now()
		_ = s.updateOrder(ctx, order)
		return order, err
	}
	order.Status = next
	order.UpdatedAt = s.clock.Now()
	records, err := s.actionRecords(ctx, order, previousStatus, action)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func bindOrderProviderConfig(provider core.Provider, order core.Order) {
	configured, ok := provider.(*ConfiguredProvider)
	if !ok {
		return
	}
	configured.BindOrderConfig(order.UpstreamOrderID)
}

func providerPolicyForOrder(ctx context.Context, provider core.Provider, order core.Order) core.ProviderPolicy {
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.LoadPolicyForOrder(ctx, order.UpstreamOrderID)
	}
	return providerPolicyForUpstreamOrder(provider, order.UpstreamOrderID)
}

func providerPolicyForUpstreamOrder(provider core.Provider, upstreamOrderID string) core.ProviderPolicy {
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.PolicyForOrder(upstreamOrderID)
	}
	return provider.Policy()
}

func shouldQueueEarlyCancelRetry(err *core.Error, policy core.ProviderPolicy) bool {
	return err != nil && err.Code == core.CodeCancelNotAllowed && err.Retryable && policy.EarlyCancelRetryAfter > 0
}

func earlyCancelRetryAt(order core.Order, policy core.ProviderPolicy, now time.Time) time.Time {
	if !order.AcquiredAt.IsZero() && policy.EarlyCancelRetryAfter > 0 {
		retryAt := order.AcquiredAt.Add(policy.EarlyCancelRetryAfter)
		if retryAt.After(now) {
			return retryAt
		}
	}
	delay := policy.PollInterval
	if delay <= 0 {
		delay = 5 * time.Second
	}
	return now.Add(delay)
}

func (s *OrderService) queueCancelRequest(ctx context.Context, order core.Order, requestID string, reason string) (core.Order, error) {
	if requestID = strings.TrimSpace(requestID); requestID == "" {
		requestID = s.ids.NewID("req_")
	}
	retryAt := s.cancelReadyAt(ctx, order)
	order = markCancelRequest(order, requestID, retryAt, s.clock.Now())
	record, err := s.events.OrderCancelRequested(ctx, order, requestID, reason)
	if err != nil {
		return order, err
	}
	if err := s.updateOrder(ctx, order, record); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func (s *OrderService) deferCancelRetry(ctx context.Context, order core.Order, requestID string, retryAt time.Time) (core.Order, error) {
	if retryAt.IsZero() {
		retryAt = s.clock.Now().Add(5 * time.Second)
	}
	if !retryAt.After(s.clock.Now()) {
		retryAt = s.clock.Now().Add(5 * time.Second)
	}
	order = markCancelRequest(order, requestID, retryAt, s.clock.Now())
	if err := s.updateOrder(ctx, order); err != nil {
		return core.Order{}, err
	}
	return order, &CancelRetryError{RetryAt: retryAt}
}

func (s *OrderService) cancelReadyAt(ctx context.Context, order core.Order) time.Time {
	now := s.clock.Now()
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return now
	}
	policy := providerPolicyForOrder(ctx, provider, order).WithDefaults()
	if policy.CancelAllowedAfter <= 0 || order.AcquiredAt.IsZero() {
		return now
	}
	readyAt := order.AcquiredAt.Add(policy.CancelAllowedAfter)
	if readyAt.After(now) {
		return readyAt
	}
	return now
}

func (s *OrderService) resolveCancelRace(ctx context.Context, order core.Order, provider core.Provider, cancelErr *core.Error) (core.Order, bool, error) {
	if cancelErr == nil {
		return order, false, nil
	}
	if cancelErr.Retryable && cancelErr.Code != core.CodeCancelNotAllowed {
		return order, false, nil
	}
	synced, handled, err := s.syncTerminalOrChargedProviderState(ctx, order, provider)
	if err != nil {
		return order, false, nil
	}
	return synced, handled, nil
}

func (s *OrderService) syncTerminalOrChargedProviderState(ctx context.Context, order core.Order, provider core.Provider) (core.Order, bool, error) {
	result, err := provider.GetStatus(ctx, order.UpstreamOrderID)
	if err != nil {
		return order, false, err
	}
	previousStatus := order.Status
	now := s.clock.Now()
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	switch result.Status {
	case core.StatusCodeReceived:
		receivedAt := result.ReceivedAt
		if receivedAt.IsZero() {
			receivedAt = now
		}
		code := core.SMSCode{Value: result.Code, MessageText: result.MessageText, ReceivedAt: receivedAt}
		order.Status = core.StatusCodeReceived
		records, err := s.statusAndCodeRecords(ctx, order, previousStatus, code)
		if err != nil {
			return order, true, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, true, err
		}
		return order, true, nil
	case core.StatusCompleted, core.StatusCanceled, core.StatusExpired, core.StatusFailed:
		order.Status = result.Status
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, true, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, true, err
		}
		return order, true, nil
	default:
		return order, false, nil
	}
}

func orderHasCode(order core.Order) bool {
	return order.Status == core.StatusCodeReceived
}

func markCancelRequest(order core.Order, requestID string, retryAt time.Time, now time.Time) core.Order {
	order.UpdatedAt = now
	order.CancelAllowedAt = retryAt
	order.LastError = nil
	return order
}

func (s *OrderService) expireIfNeeded(ctx context.Context, order *core.Order) error {
	now := s.clock.Now()
	if !order.IsExpired(now) {
		return nil
	}
	previousStatus := order.Status
	order.Status = core.StatusExpired
	order.UpdatedAt = now
	records, err := s.statusChangedRecords(ctx, *order, previousStatus)
	if err != nil {
		return err
	}
	if err := s.updateOrder(ctx, *order, records...); err != nil {
		return err
	}
	return core.NewError(core.CodeOrderExpired, "order expired", false)
}

func (s *OrderService) orderAcquiredRecords(ctx context.Context, order core.Order, reason string) ([]eventoutbox.Record, error) {
	acquired, err := s.events.OrderAcquired(ctx, order)
	if err != nil {
		return nil, err
	}
	poll, err := s.events.OrderPollRequested(ctx, order, reason)
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(acquired, poll), nil
}

func (s *OrderService) statusAndCodeRecords(ctx context.Context, order core.Order, previous core.OrderStatus, code core.SMSCode) ([]eventoutbox.Record, error) {
	records, err := s.statusChangedRecords(ctx, order, previous)
	if err != nil {
		return nil, err
	}
	codeRecord, err := s.events.CodeReceived(ctx, order, code)
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(append(records, codeRecord)...), nil
}

func (s *OrderService) actionRecords(ctx context.Context, order core.Order, previous core.OrderStatus, action core.ProviderAction) ([]eventoutbox.Record, error) {
	records, err := s.statusChangedRecords(ctx, order, previous)
	if err != nil {
		return nil, err
	}
	if action != core.ActionMarkMessageSent && action != core.ActionRequestAdditional {
		return records, nil
	}
	poll, err := s.events.OrderPollRequested(ctx, order, string(action))
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(append(records, poll)...), nil
}

func (s *OrderService) statusChangedRecords(ctx context.Context, order core.Order, previous core.OrderStatus) ([]eventoutbox.Record, error) {
	if previous == order.Status {
		return nil, nil
	}
	record, err := s.events.OrderStatusChanged(ctx, order, previous)
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(record), nil
}

func nonEmptyRecords(records ...eventoutbox.Record) []eventoutbox.Record {
	out := make([]eventoutbox.Record, 0, len(records))
	for _, record := range records {
		if strings.TrimSpace(record.EventID) == "" {
			continue
		}
		out = append(out, record)
	}
	return out
}

func (s *OrderService) provider(key string) (core.Provider, error) {
	provider, ok := s.providers[key]
	if !ok {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider not registered", false)
	}
	return provider, nil
}

func (s *OrderService) routePolicy(ctx context.Context, route core.Route) core.ProviderPolicy {
	provider, err := s.provider(route.ProviderKey)
	if err != nil {
		return core.ProviderPolicy{}.WithDefaults()
	}
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.LoadPolicyForOrder(ctx, "").WithDefaults()
	}
	return provider.Policy().WithDefaults()
}

func withRouteTargetDefaults(target core.Target, route core.Route) core.Target {
	if target.ApplicationKey == "" {
		target.ApplicationKey = route.ApplicationKey
	}
	if target.CountryISO2 == "" {
		target.CountryISO2 = route.CountryISO2
	}
	if target.CountryCallingCode == "" {
		target.CountryCallingCode = route.CountryCallingCode
	}
	return target
}

func routeFromOrder(order core.Order) core.Route {
	return core.Route{
		ProviderKey:        order.ProviderKey,
		ApplicationKey:     order.Target.ApplicationKey,
		UpstreamServiceKey: order.Target.ApplicationKey,
		CountryISO2:        order.Target.CountryISO2,
		CountryCallingCode: order.Target.CountryCallingCode,
	}
}

func orderRequestExpiresAt(now time.Time, policy core.ProviderPolicy, lease time.Duration) time.Time {
	policy = policy.WithDefaults()
	ttl := policy.OrderTTL
	if lease > 0 && lease < ttl {
		ttl = lease
	}
	if ttl <= 0 {
		return time.Time{}
	}
	return now.Add(ttl)
}

func remainingLease(now time.Time, expiresAt time.Time) time.Duration {
	if expiresAt.IsZero() {
		return 0
	}
	if !expiresAt.After(now) {
		return 0
	}
	return expiresAt.Sub(now)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func asCoreError(err error) *core.Error {
	if err == nil {
		return nil
	}
	if smsErr, ok := err.(*core.Error); ok {
		return smsErr
	}
	return core.NewError(core.CodeInternal, err.Error(), false)
}
