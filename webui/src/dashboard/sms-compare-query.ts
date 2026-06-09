import type { SmsPriceOfferQuery } from './sms-api';
import type { OfferSort } from './sms-compare-data';

export type CompareQuery = {
  serviceText: string;
  countryISO2: string;
  countryCallingCode: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
};

export function canUseCatalog(providerKeys: string[]) {
  return providerKeys.length > 0;
}

export function canSearch(providerKeys: string[], serviceText: string) {
  return canUseCatalog(providerKeys) && serviceText.trim().length > 0;
}

export function compareQueryFromSearch(params: URLSearchParams): CompareQuery {
  return {
    serviceText: params.get('q') || params.get('application_key') || '',
    countryISO2: params.get('country_iso2') || '',
    countryCallingCode: params.get('country_calling_code') || '',
    providerKeys: [],
    minAvailable: Math.max(0, numberInputValue(params.get('min_available') || '1')),
    sort: offerSort(params.get('sort'))
  };
}

export function compareQuerySearchParams(query: CompareQuery) {
  const params = new URLSearchParams();
  if (query.serviceText) params.set('q', query.serviceText);
  if (query.countryISO2) params.set('country_iso2', query.countryISO2);
  if (query.countryCallingCode) params.set('country_calling_code', query.countryCallingCode);
  params.set('min_available', String(query.minAvailable));
  params.set('sort', query.sort);
  return params;
}

export function smsPriceOfferQuery(query: CompareQuery): SmsPriceOfferQuery {
  return {
    providerKeys: query.providerKeys,
    applicationKey: query.serviceText || undefined,
    countryISO2: query.countryISO2 || undefined,
    countryCallingCode: query.countryCallingCode || undefined
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
