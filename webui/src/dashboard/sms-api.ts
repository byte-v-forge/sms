import { api } from '@/dashboard/module-kit';
import type {
  CancelProviderActivationResponse,
  DeleteProviderConfigResponse,
  GetProviderBalanceResponse,
  ListProviderPluginsResponse,
  ListRouteOptionsResponse,
  ListActivationsResponse,
  ListProviderConfigsResponse,
  ListRouteProfilesResponse,
  SmsProviderConfig,
  SmsRouteProfile,
  UpsertRouteProfileResponse,
  DeleteRouteProfileResponse,
  UpsertProviderConfigResponse
} from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';

export const smsKeys = {
  plugins: ['sms', 'provider-plugins'] as const,
  configs: ['sms', 'provider-configs'] as const,
  profiles: ['sms', 'route-profiles'] as const,
  activations: ['sms', 'activations'] as const,
  balance: (id: string) => ['sms', 'balance', id] as const,
  routeOptions: (providerKey: string) => ['sms', 'route-options', providerKey] as const
};

export function listSmsProviderPlugins() {
  return api<ListProviderPluginsResponse>('/api/sms/provider-plugins');
}

export function listSmsProviderConfigs() {
  return api<ListProviderConfigsResponse>('/api/sms/provider-configs?include_disabled=true');
}

export function saveSmsProviderConfig(config: SmsProviderConfig) {
  return api<UpsertProviderConfigResponse>('/api/sms/provider-configs', {
    method: 'POST',
    body: JSON.stringify({ config })
  });
}

export function deleteSmsProviderConfig(id: string) {
  return api<DeleteProviderConfigResponse>(`/api/sms/provider-configs/${encodeURIComponent(id)}`, { method: 'DELETE' });
}

export function listSmsRouteProfiles() {
  return api<ListRouteProfilesResponse>('/api/sms/route-profiles?include_disabled=true');
}

export function saveSmsRouteProfile(profile: SmsRouteProfile) {
  return api<UpsertRouteProfileResponse>('/api/sms/route-profiles', {
    method: 'POST',
    body: JSON.stringify({ profile })
  });
}

export function deleteSmsRouteProfile(profileKey: string) {
  return api<DeleteRouteProfileResponse>(`/api/sms/route-profiles/${encodeURIComponent(profileKey)}`, { method: 'DELETE' });
}

export function listSmsRouteOptions(providerKey: string) {
  return api<ListRouteOptionsResponse>(`/api/sms/route-options?provider_key=${encodeURIComponent(providerKey)}`);
}

export function getSmsProviderBalance(id: string) {
  return api<GetProviderBalanceResponse>(`/api/sms/provider-configs/${encodeURIComponent(id)}/balance`);
}

export function listSmsActivations() {
  return api<ListActivationsResponse>('/api/sms/activations?limit=200');
}

export function cancelSmsActivation(id: string) {
  return api<CancelProviderActivationResponse>(`/api/sms/activations/${encodeURIComponent(id)}/cancel`, {
    method: 'POST',
    body: JSON.stringify({})
  });
}
