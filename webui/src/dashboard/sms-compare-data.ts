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

export function filterAndSortOffers(offers: SmsPriceOffer[], providerKeys: string[], searchText: string, minAvailable: number, sort: OfferSort) {
  const providerFilter = new Set(providerKeys);
  const tokens = searchTokens(searchText);
  return offers
    .filter((offer) => providerFilter.size === 0 || providerFilter.has(offer.provider_key))
    .filter((offer) => offerMatchesSearch(offer, tokens))
    .filter((offer) => offer.available_count >= minAvailable)
    .sort((left, right) => compareOffers(left, right, sort));
}

export function bestOffer(offers: SmsPriceOffer[]) {
  return [...offers].sort((left, right) => compareOffers(left, right, 'price'))[0];
}

export function offerRowKey(offer: SmsPriceOffer) {
  return offer.offer_ref?.offer_id || [offer.provider_key, offer.application_key, offer.country_iso2, offer.country_calling_code, offer.price?.amount_decimal].join(':');
}

function compareOffers(left: SmsPriceOffer, right: SmsPriceOffer, sort: OfferSort) {
  if (sort === 'available') return right.available_count - left.available_count || compareMoney(left, right);
  if (sort === 'provider') return left.provider_display_name.localeCompare(right.provider_display_name) || compareMoney(left, right);
  return compareMoney(left, right) || right.available_count - left.available_count;
}

function compareMoney(left: SmsPriceOffer, right: SmsPriceOffer) {
  return moneyAmount(left) - moneyAmount(right);
}

function moneyAmount(offer: SmsPriceOffer) {
  const amount = Number(offer.price?.amount_decimal || Number.POSITIVE_INFINITY);
  return Number.isFinite(amount) ? amount : Number.POSITIVE_INFINITY;
}

function searchTokens(text: string) {
  return text.toLowerCase().split(/\s+/).map((token) => token.trim()).filter(Boolean);
}

function offerMatchesSearch(offer: SmsPriceOffer, tokens: string[]) {
  if (tokens.length === 0) return true;
  const haystack = [
    offer.provider_key,
    offer.provider_display_name,
    offer.application_key,
    offer.application_name,
    offer.country_iso2,
    offer.country_name,
    offer.country_calling_code && `+${offer.country_calling_code}`,
    offer.offer_ref?.offer_id,
    offer.offer_ref?.route_ref?.upstream_service_key,
    offer.offer_ref?.route_ref?.provider_country_id,
    offer.offer_ref?.route_ref?.upstream_provider_id
  ].filter(Boolean).join(' ').toLowerCase();
  return tokens.every((token) => haystack.includes(token));
}
