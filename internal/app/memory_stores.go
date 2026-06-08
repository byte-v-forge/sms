package app

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TTLStringStore interface {
	DefaultTTL() time.Duration
	Load(ctx context.Context, key string) (string, bool, error)
	SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error
}

type MemoryProviderConfigStore struct {
	mu        sync.RWMutex
	providers *providerspi.Registry
	configs   map[string]*smsinternalv1.SmsProviderConfig
	clock     core.Clock
}

func NewMemoryProviderConfigStore(providers *providerspi.Registry, clock core.Clock) *MemoryProviderConfigStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryProviderConfigStore{providers: providers, configs: map[string]*smsinternalv1.SmsProviderConfig{}, clock: clock}
}

func (s *MemoryProviderConfigStore) UpsertProviderConfig(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	config, err := s.normalizeForSave(input)
	if err != nil {
		return nil, err
	}
	now := timestamppb.New(s.clock.Now())
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing := s.configs[config.GetProviderKey()]; existing != nil {
		config.CreatedAt = cloneTimestamp(existing.GetCreatedAt())
		if strings.TrimSpace(config.GetCredentialSecret()) == "" {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
	} else {
		config.CreatedAt = cloneTimestamp(now)
	}
	config.UpdatedAt = cloneTimestamp(now)
	config.CredentialSecretSet = strings.TrimSpace(config.GetCredentialSecret()) != ""
	s.configs[config.GetProviderKey()] = cloneProviderConfig(config)
	return cloneProviderConfig(config), nil
}

func (s *MemoryProviderConfigStore) GetProviderConfig(ctx context.Context, providerKey string) (*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	providerKey = normalizeProviderKey(providerKey)
	if providerKey == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	config := s.configs[providerKey]
	if config == nil {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	return cloneProviderConfig(config), nil
}

func (s *MemoryProviderConfigStore) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	providerKey = normalizeProviderKey(providerKey)
	s.mu.RLock()
	defer s.mu.RUnlock()
	configs := make([]*smsinternalv1.SmsProviderConfig, 0, len(s.configs))
	for key, config := range s.configs {
		if providerKey != "" && key != providerKey {
			continue
		}
		if !includeDisabled && !config.GetEnabled() {
			continue
		}
		configs = append(configs, cloneProviderConfig(config))
	}
	sort.Slice(configs, func(i, j int) bool { return configs[i].GetProviderKey() < configs[j].GetProviderKey() })
	return configs, nil
}

func (s *MemoryProviderConfigStore) DeleteProviderConfig(ctx context.Context, providerKey string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	providerKey = normalizeProviderKey(providerKey)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.configs[providerKey] == nil {
		return core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	delete(s.configs, providerKey)
	return nil
}

func (s *MemoryProviderConfigStore) GetEnabledProviderConfig(ctx context.Context, providerKey string, _ core.Target) (*smsinternalv1.SmsProviderConfig, error) {
	configs, err := s.ListProviderConfigs(ctx, false, providerKey)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, core.NewError(core.CodeRouteNotFound, "no enabled sms provider config", false)
	}
	return configs[0], nil
}

func (s *MemoryProviderConfigStore) normalizeForSave(input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config := cloneProviderConfig(input)
	config.ProviderKey = normalizeProviderKey(config.GetProviderKey())
	if config.GetProviderKey() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	if s.providers != nil && !s.providers.Supports(config.GetProviderKey()) {
		return nil, core.NewError(core.CodeUnsupportedOperation, "unsupported sms provider", false)
	}
	config.CredentialSecret = strings.TrimSpace(config.GetCredentialSecret())
	if config.GetCredentialSecret() == "" {
		s.mu.RLock()
		existing := s.configs[config.GetProviderKey()]
		if existing != nil {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
		s.mu.RUnlock()
	}
	if config.GetEnabled() && config.GetCredentialSecret() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "credential_secret is required for enabled sms provider", false)
	}
	return config, nil
}

type MemoryOrderStore struct {
	mu     sync.RWMutex
	orders map[string]core.Order
	codes  map[string][]core.OrderCode
}

func NewMemoryOrderStore() *MemoryOrderStore {
	return &MemoryOrderStore{orders: map[string]core.Order{}, codes: map[string][]core.OrderCode{}}
}

func (s *MemoryOrderStore) Save(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *MemoryOrderStore) Update(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *MemoryOrderStore) RecordCode(ctx context.Context, order core.Order, code core.SMSCode, _ ...eventoutbox.Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = cloneOrder(order)
	s.codes[order.ID] = append(s.codes[order.ID], core.OrderCode{OrderID: order.ID, Code: cloneSMSCode(code)})
	return nil
}

func (s *MemoryOrderStore) CodeSecretExists(ctx context.Context, orderID string, secretID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, code := range s.codes[orderID] {
		if code.Code.SecretRef.GetSecretId() == secretID {
			return true, nil
		}
	}
	return false, nil
}

func (s *MemoryOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	if err := ctx.Err(); err != nil {
		return core.Order{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderID]
	if !ok {
		return core.Order{}, core.NewError(core.CodeOrderNotFound, "order not found", false)
	}
	return cloneOrder(order), nil
}

func (s *MemoryOrderStore) List(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]core.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if !includeFinal && order.Status.IsFinal() {
			continue
		}
		orders = append(orders, cloneOrder(order))
	}
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].UpdatedAt.Equal(orders[j].UpdatedAt) {
			return orders[i].ID < orders[j].ID
		}
		return orders[i].UpdatedAt.After(orders[j].UpdatedAt)
	})
	if len(orders) > limit {
		return orders[:limit], nil
	}
	return orders, nil
}

func (s *MemoryOrderStore) ListCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if limitPerOrder <= 0 {
		limitPerOrder = 10
	}
	if limitPerOrder > 50 {
		limitPerOrder = 50
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []core.OrderCode{}
	for _, orderID := range orderIDs {
		codes := append([]core.OrderCode{}, s.codes[orderID]...)
		sort.Slice(codes, func(i, j int) bool {
			return codes[i].Code.ReceivedAt.After(codes[j].Code.ReceivedAt)
		})
		if len(codes) > limitPerOrder {
			codes = codes[:limitPerOrder]
		}
		for _, code := range codes {
			out = append(out, core.OrderCode{OrderID: code.OrderID, Code: cloneSMSCode(code.Code)})
		}
	}
	return out, nil
}

func (s *MemoryOrderStore) save(ctx context.Context, order core.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = cloneOrder(order)
	return nil
}

type MemoryTTLStringStore struct {
	mu       sync.RWMutex
	values   map[string]memoryTTLValue
	ttl      time.Duration
	clock    core.Clock
	keyspace string
}

type memoryTTLValue struct {
	value     string
	expiresAt time.Time
}

func NewMemoryTTLStringStore(prefix string, ttl time.Duration, clock core.Clock) *MemoryTTLStringStore {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryTTLStringStore{values: map[string]memoryTTLValue{}, ttl: ttl, clock: clock, keyspace: strings.Trim(strings.TrimSpace(prefix), ":")}
}

func (s *MemoryTTLStringStore) DefaultTTL() time.Duration {
	if s == nil {
		return 0
	}
	return s.ttl
}

func (s *MemoryTTLStringStore) Load(ctx context.Context, key string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}
	normalized := s.key(key)
	if normalized == "" {
		return "", false, nil
	}
	s.mu.RLock()
	item, ok := s.values[normalized]
	s.mu.RUnlock()
	if !ok {
		return "", false, nil
	}
	if !item.expiresAt.IsZero() && !s.clock.Now().Before(item.expiresAt) {
		s.mu.Lock()
		delete(s.values, normalized)
		s.mu.Unlock()
		return "", false, nil
	}
	return item.value, true, nil
}

func (s *MemoryTTLStringStore) SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	normalized := s.key(key)
	if normalized == "" {
		return core.NewError(core.CodeValidationFailed, "memory string store key is required", false)
	}
	if ttl <= 0 {
		ttl = s.ttl
	}
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = s.clock.Now().Add(ttl)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[normalized] = memoryTTLValue{value: value, expiresAt: expiresAt}
	return nil
}

func (s *MemoryTTLStringStore) key(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if s.keyspace == "" {
		return value
	}
	return s.keyspace + ":" + value
}

func cloneOrder(order core.Order) core.Order {
	order.LastError = cloneCoreError(order.LastError)
	return order
}

func cloneCoreError(err *core.Error) *core.Error {
	if err == nil {
		return nil
	}
	out := *err
	return &out
}

func cloneSMSCode(code core.SMSCode) core.SMSCode {
	code.SecretRef = cloneSecretRef(code.SecretRef)
	return code
}

func cloneSecretRef(ref *commonv1.SecretRef) *commonv1.SecretRef {
	if ref == nil {
		return nil
	}
	return proto.Clone(ref).(*commonv1.SecretRef)
}

func cloneTimestamp(ts *timestamppb.Timestamp) *timestamppb.Timestamp {
	if ts == nil {
		return nil
	}
	return timestamppb.New(ts.AsTime())
}
