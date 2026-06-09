import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';

export type OfferSort = 'price' | 'available' | 'provider';

export type ProviderChoice = {
  providerKey: string;
  displayName: string;
  enabled: boolean;
  configured: boolean;
};

export function providerChoices(options: SmsProviderPluginDescriptor[], configs: SmsProviderConfig[]): ProviderChoice[] {
  const configByKey = new Map(configs.map((config) => [config.provider_key, config]));
  const choices = options.map((option) => {
    const config = configByKey.get(option.provider_key);
    return {
      providerKey: option.provider_key,
      displayName: option.display_name || option.provider_key,
      enabled: !!config?.enabled,
      configured: !!config?.provider_key
    };
  });
  const optionKeys = new Set(options.map((option) => option.provider_key));
  for (const config of configs) {
    if (!optionKeys.has(config.provider_key)) {
      choices.push({ providerKey: config.provider_key, displayName: config.provider_key, enabled: !!config.enabled, configured: true });
    }
  }
  return choices.sort((left, right) => left.displayName.localeCompare(right.displayName));
}

export function enabledProviderKeys(choices: ProviderChoice[]) {
  return choices.filter((choice) => choice.enabled).map((choice) => choice.providerKey);
}

export function filterAndSortOffers(offers: SmsPriceOffer[], providerKeys: string[], countryISO2: string, countryCallingCode: string, minAvailable: number, sort: OfferSort) {
  const providerFilter = new Set(providerKeys);
  return offers
    .filter((offer) => providerFilter.size === 0 || providerFilter.has(offer.provider_key))
    .filter((offer) => countryMatches(offer, countryISO2, countryCallingCode))
    .filter((offer) => availableCount(offer) >= minAvailable)
    .sort((left, right) => compareOffers(left, right, sort));
}

export function bestOffer(offers: SmsPriceOffer[]) {
  return [...offers].sort((left, right) => compareOffers(left, right, 'price'))[0];
}

export function availableCount(offer: SmsPriceOffer) {
  return Number(offer.available_count || 0);
}

export function offerRowKey(offer: SmsPriceOffer) {
  return offer.offer_ref?.offer_id || [offer.provider_key, offer.application_key, offer.country_iso2, offer.country_calling_code, offer.price?.amount_decimal].join(':');
}

function compareOffers(left: SmsPriceOffer, right: SmsPriceOffer, sort: OfferSort) {
  if (sort === 'available') return availableCount(right) - availableCount(left) || compareMoney(left, right);
  if (sort === 'provider') return left.provider_display_name.localeCompare(right.provider_display_name) || compareMoney(left, right);
  return compareMoney(left, right) || availableCount(right) - availableCount(left);
}

function compareMoney(left: SmsPriceOffer, right: SmsPriceOffer) {
  return moneyAmount(left) - moneyAmount(right);
}

function moneyAmount(offer: SmsPriceOffer) {
  const amount = Number(offer.price?.amount_decimal || Number.POSITIVE_INFINITY);
  return Number.isFinite(amount) ? amount : Number.POSITIVE_INFINITY;
}

function countryMatches(offer: SmsPriceOffer, countryISO2: string, countryCallingCode: string) {
  if (countryISO2 && offer.country_iso2.toLowerCase() !== countryISO2.toLowerCase()) return false;
  return !(countryCallingCode && offer.country_calling_code !== countryCallingCode);
}
