import type { ListSmsPriceOffersResponse, SmsError, SmsProviderLookupError } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { SmsErrorCode } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { listSmsApplications, listSmsPriceOffers } from './sms-api';
import { applicationChoices, matchApplicationChoice, providerApplicationKey, type ApplicationChoice } from './sms-compare-options';
import type { CompareQuery } from './sms-compare-query';

export async function listSmsProviderOfferSearch(query: CompareQuery, application: ApplicationChoice | undefined) {
  const responses = await Promise.all(query.providerKeys.map((providerKey) => listProviderOfferSearch(providerKey, query, application)));
  return mergeSmsPriceOfferResponses(responses);
}

async function listProviderOfferSearch(providerKey: string, query: CompareQuery, application: ApplicationChoice | undefined) {
  try {
    const applicationKey = await providerSearchApplicationKey(providerKey, query.serviceText, application);
    return await listSmsPriceOffers({
      providerKeys: [providerKey],
      applicationKey,
      countryISO2: query.countryISO2 || undefined,
      countryCallingCode: query.countryCallingCode || undefined
    });
  } catch (error) {
    return { offers: [], provider_errors: [providerLookupError(providerKey, error)], error: undefined };
  }
}

async function providerSearchApplicationKey(providerKey: string, serviceText: string, application: ApplicationChoice | undefined) {
  const selectedKey = application?.providerApplicationKeys[providerKey];
  if (selectedKey) return selectedKey;
  const response = await listSmsApplications({ providerKeys: [providerKey], searchText: serviceText });
  const choices = applicationChoices([{ providerKey, applications: response.applications || [], error: response.error }], []);
  const choice = matchApplicationChoice(serviceText, choices) || singleChoice(choices);
  return providerApplicationKey(choice, providerKey, serviceText);
}

function singleChoice(choices: ApplicationChoice[]) {
  return choices.length === 1 ? choices[0] : undefined;
}

function mergeSmsPriceOfferResponses(responses: ListSmsPriceOffersResponse[]): ListSmsPriceOffersResponse {
  const offers = responses.flatMap((response) => response.offers || []);
  const providerErrors = responses.flatMap((response) => response.provider_errors || []);
  const error = offers.length === 0 ? responses.find((response) => response.error)?.error : undefined;
  return { offers, provider_errors: providerErrors, error };
}

function providerLookupError(providerKey: string, error: unknown): SmsProviderLookupError {
  return { provider_key: providerKey, provider_display_name: providerKey, error: internalSmsError(error) };
}

function internalSmsError(error: unknown): SmsError {
  return { code: SmsErrorCode.SMS_ERROR_CODE_INTERNAL, message: error instanceof Error ? error.message : String(error), retryable: true };
}
