import { api } from '@byte-v-forge/common-ui';
import type {
  CancelProviderOrderResponse,
  DeleteProviderConfigResponse,
  GetProviderBalanceResponse,
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

export const smsKeys = {
  settingsProviders: ['sms', 'settings', 'providers'] as const,
  orders: ['sms', 'orders'] as const,
  balance: (providerKey: string) => ['sms', 'balance', providerKey] as const
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
  return api<ListOrdersResponse>('/api/sms/orders?limit=200');
}

export function cancelSmsOrder(id: string) {
  return api<CancelProviderOrderResponse>(`/api/sms/orders/${encodeURIComponent(id)}/cancel`, { method: 'POST', body: JSON.stringify({}) });
}
