import type { SmsPriceOfferQuery } from './sms-api';
import type { OfferSort } from './sms-compare-data';

export type CompareQuery = {
  searchText: string;
  applicationKey: string;
  countryISO2: string;
  countryCallingCode: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
};

export function canSearch(providerKeys: string[]) {
  return providerKeys.length > 0;
}

export function compareQueryFromSearch(params: URLSearchParams): CompareQuery {
  return {
    searchText: params.get('q') || '',
    applicationKey: params.get('application_key') || '',
    countryISO2: params.get('country_iso2') || '',
    countryCallingCode: params.get('country_calling_code') || '',
    providerKeys: params.getAll('provider_key').filter(Boolean),
    minAvailable: Math.max(0, numberInputValue(params.get('min_available') || '1')),
    sort: offerSort(params.get('sort'))
  };
}

export function compareQuerySearchParams(query: CompareQuery) {
  const params = new URLSearchParams();
  if (query.searchText) params.set('q', query.searchText);
  if (query.applicationKey) params.set('application_key', query.applicationKey);
  if (query.countryISO2) params.set('country_iso2', query.countryISO2);
  if (query.countryCallingCode) params.set('country_calling_code', query.countryCallingCode);
  for (const providerKey of query.providerKeys) params.append('provider_key', providerKey);
  params.set('min_available', String(query.minAvailable));
  params.set('sort', query.sort);
  return params;
}

export function smsPriceOfferQuery(query: CompareQuery): SmsPriceOfferQuery {
  return {
    providerKeys: query.providerKeys
  };
}

export function numberInputValue(value: string | number | null) {
  const number = Number(value || 0);
  return Number.isFinite(number) ? number : 0;
}

function offerSort(value: string | null): OfferSort {
  if (value === 'available' || value === 'provider') return value;
  return 'price';
}
