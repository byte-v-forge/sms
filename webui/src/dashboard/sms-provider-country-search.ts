import type { ListSmsCountriesResponse, SmsCountry, SmsError } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { SmsErrorCode } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { listSmsCountries } from './sms-api';
import type { ApplicationChoice } from './sms-compare-options';
import type { CompareQuery } from './sms-compare-query';
import { resolveSmsProviderApplicationKey } from './sms-provider-application-resolver';

export async function listSmsProviderCountrySearch(query: CompareQuery, application: ApplicationChoice | undefined) {
  const responses = await Promise.all(query.providerKeys.map((providerKey) => listProviderCountrySearch(providerKey, query, application)));
  return mergeSmsCountryResponses(responses);
}

async function listProviderCountrySearch(providerKey: string, query: CompareQuery, application: ApplicationChoice | undefined) {
  try {
    const applicationKey = await resolveSmsProviderApplicationKey(providerKey, query.serviceText, application);
    return await listSmsCountries({ providerKeys: [providerKey], applicationKey });
  } catch (error) {
    return { countries: [], error: internalSmsError(providerKey, error) };
  }
}

function mergeSmsCountryResponses(responses: ListSmsCountriesResponse[]): ListSmsCountriesResponse {
  const countries = uniqueCountries(responses.flatMap((response) => response.countries || []));
  const error = countries.length === 0 ? responses.find((response) => response.error)?.error : undefined;
  return { countries, error };
}

function uniqueCountries(countries: SmsCountry[]) {
  const items = new Map<string, SmsCountry>();
  for (const country of countries) {
    const key = [country.country_iso2.toUpperCase(), country.country_calling_code].join('|');
    if (key === '|') continue;
    const current = items.get(key);
    items.set(key, current && current.name.length >= country.name.length ? current : country);
  }
  return [...items.values()].sort((left, right) => (left.name || left.country_iso2).localeCompare(right.name || right.country_iso2));
}

function internalSmsError(providerKey: string, error: unknown): SmsError {
  const message = error instanceof Error ? error.message : String(error);
  return { code: SmsErrorCode.SMS_ERROR_CODE_INTERNAL, message: `${providerKey}: ${message}`, retryable: true };
}
