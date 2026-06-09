import type { SmsProviderPriceOfferQuery } from './sms-api';
import type { ApplicationChoice } from './sms-compare-options';
import { providerApplicationKey } from './sms-compare-options';
import type { CompareQuery } from './sms-compare-query';

export function providerOfferQueries(query: CompareQuery, application: ApplicationChoice | undefined): SmsProviderPriceOfferQuery[] {
  const fallback = (application?.displayName || query.serviceText).trim();
  if (!fallback) return [];
  return query.providerKeys.map((providerKey) => ({
    providerKey,
    applicationKey: providerApplicationKey(application, providerKey, fallback),
    countryISO2: query.countryISO2 || undefined,
    countryCallingCode: query.countryCallingCode || undefined
  }));
}
