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
  SmsProviderConfig,
  SmsProviderPluginDescriptor
} from '../proto/byte/v/forge/sms/internal/v1/sms_internal';

export type SmsProviderOption = Pick<SmsProviderPluginDescriptor, 'provider_key' | 'display_name'>;
export type SmsProviderSetting = Pick<SmsProviderConfig, 'provider_key' | 'enabled'> & {
  api_key_set?: boolean;
};
export type ListSmsProviderSettingsResponse = {
  provider_options?: SmsProviderOption[];
  providers?: SmsProviderSetting[];
};
export type SaveSmsProviderSettingRequest = {
  provider_key: string;
  enabled?: boolean;
  api_key?: string;
};
export type SaveSmsProviderSettingResponse = {
  provider?: SmsProviderSetting;
};

export type SmsPriceOfferQuery = {
  applicationKey: string;
  countryISO2: string;
  countryCallingCode: string;
  providerKeys: string[];
  minAvailable: number;
};

export const smsKeys = {
  settingsProviders: ['sms', 'settings', 'providers'] as const,
  orders: ['sms', 'orders'] as const,
  orderCodesRoot: ['sms', 'order-codes'] as const,
  orderCodes: (orderIds: string[]) => ['sms', 'order-codes', orderIds.join(',')] as const,
  balance: (providerKey: string) => ['sms', 'balance', providerKey] as const,
  priceOffers: (query?: SmsPriceOfferQuery) => ['sms', 'price-offers', query] as const
};

export function listSmsProviderSettings() {
  return api<ListSmsProviderSettingsResponse>('/api/sms/settings/providers');
}

export function saveSmsProviderSetting(input: SaveSmsProviderSettingRequest) {
  return api<SaveSmsProviderSettingResponse>('/api/sms/settings/providers', { method: 'POST', body: JSON.stringify(input) });
}

export function deleteSmsProviderSetting(providerKey: string) {
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
  if (query.providerKeys.length === 1) params.set('provider_key', query.providerKeys[0]);
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
