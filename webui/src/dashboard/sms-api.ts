import { api } from '../ui';
import type {
  AcquireNumberRequest,
  AcquireNumberResponse,
  ListSmsPriceOffersResponse,
  SmsPriceOffer
} from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type {
  CancelProviderOrderResponse,
  DeleteProviderConfigResponse,
  GetProviderBalanceResponse,
  ListOrderCodesResponse,
  ListOrdersResponse,
  ListProviderConfigsResponse,
  ListProviderPluginsResponse,
  SmsProviderConfig,
  UpsertProviderConfigRequest,
  UpsertProviderConfigResponse
} from '../proto/byte/v/forge/sms/internal/v1/sms_internal';

export type SmsPriceOfferQuery = {
  applicationKey: string;
  countryISO2: string;
  countryCallingCode: string;
  providerKeys: string[];
  minAvailable: number;
};

export const smsKeys = {
  providerPlugins: ['sms', 'settings', 'provider-plugins'] as const,
  providerConfigs: ['sms', 'settings', 'provider-configs'] as const,
  orders: ['sms', 'orders'] as const,
  orderCodesRoot: ['sms', 'order-codes'] as const,
  orderCodes: (orderIds: string[]) => ['sms', 'order-codes', orderIds.join(',')] as const,
  balance: (providerKey: string) => ['sms', 'balance', providerKey] as const,
  priceOffers: (query?: SmsPriceOfferQuery) => ['sms', 'price-offers', query] as const
};

export function listSmsProviderPlugins() {
  return api<ListProviderPluginsResponse>('/api/sms/settings/provider-plugins');
}

export function listSmsProviderConfigs() {
  return api<ListProviderConfigsResponse>('/api/sms/settings/providers');
}

export function saveSmsProviderConfig(config: SmsProviderConfig) {
  const request: UpsertProviderConfigRequest = { config };
  return api<UpsertProviderConfigResponse>('/api/sms/settings/providers', { method: 'POST', body: JSON.stringify(request) });
}

export function deleteSmsProviderConfig(providerKey: string) {
  return api<DeleteProviderConfigResponse>(`/api/sms/settings/providers/${encodeURIComponent(providerKey)}`, { method: 'DELETE' });
}

export function getSmsProviderBalance(providerKey: string) {
  return api<GetProviderBalanceResponse>(`/api/sms/settings/providers/${encodeURIComponent(providerKey)}/balance`);
}

export function listSmsOrders() {
  return api<ListOrdersResponse>('/api/sms/orders?include_final=true&limit=200');
}

export function listSmsOrderCodes(orderIds: string[], limitPerOrder = 5) {
  const params = new URLSearchParams({ limit_per_order: String(limitPerOrder) });
  orderIds.forEach((id) => params.append('order_id', id));
  return api<ListOrderCodesResponse>(`/api/sms/order-codes?${params.toString()}`);
}

export function cancelSmsOrder(id: string) {
  return api<CancelProviderOrderResponse>(`/api/sms/orders/${encodeURIComponent(id)}/cancel`, { method: 'POST', body: JSON.stringify({}) });
}

export function listSmsPriceOffers(query: SmsPriceOfferQuery) {
  const params = new URLSearchParams();
  params.set('application_key', query.applicationKey);
  if (query.countryISO2) params.set('country_iso2', query.countryISO2);
  if (query.countryCallingCode) params.set('country_calling_code', query.countryCallingCode);
  query.providerKeys.forEach((providerKey) => params.append('provider_key', providerKey));
  return api<ListSmsPriceOffersResponse>(`/api/sms/price-offers?${params.toString()}`);
}

export function acquireSmsFromOffer(offer: SmsPriceOffer) {
  const request: AcquireNumberRequest = {
    request_id: randomRequestID(),
    lease_duration: undefined,
    acquire_params: {
      offer_ref: offer.offer_ref,
      application_key: offer.application_key,
      country_iso2: offer.country_iso2,
      country_calling_code: offer.country_calling_code,
      min_available_count: 1,
      route_failure_policy: undefined,
      max_price: offer.price,
      min_price: undefined
    }
  };
  return api<AcquireNumberResponse>('/api/sms/orders/acquire?wait_seconds=60', {
    method: 'POST',
    body: JSON.stringify(request)
  });
}

function randomRequestID() {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) return crypto.randomUUID();
  return `sms-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}
