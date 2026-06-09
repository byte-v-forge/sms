import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';

export type OfferSort = 'price' | 'available' | 'provider';

export type ProviderChoice = {
  providerKey: string;
  displayName: string;
  enabled: boolean;
  configured: boolean;
};

export type ApplicationChoice = {
  applicationKey: string;
  displayName: string;
  offerCount: number;
};

export type CountryChoice = {
  countryISO2: string;
  countryName: string;
  countryCallingCode: string;
  offerCount: number;
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

export function applicationChoices(offers: SmsPriceOffer[], providerKeys: string[]) {
  const items = new Map<string, ApplicationChoice>();
  for (const offer of providerOffers(offers, providerKeys)) {
    if (!offer.application_key) continue;
    const current = items.get(offer.application_key);
    items.set(offer.application_key, {
      applicationKey: offer.application_key,
      displayName: current?.displayName || offer.application_name || offer.application_key,
      offerCount: (current?.offerCount || 0) + 1
    });
  }
  return [...items.values()].sort((left, right) => applicationChoiceLabel(left).localeCompare(applicationChoiceLabel(right)));
}

export function countryChoices(offers: SmsPriceOffer[], providerKeys: string[], applicationKey: string, searchText: string) {
  const items = new Map<string, CountryChoice>();
  for (const offer of providerOffers(offers, providerKeys).filter((item) => applicationMatches(item, applicationKey, searchText))) {
    const key = countryChoiceValue({ countryISO2: offer.country_iso2, countryName: offer.country_name, countryCallingCode: offer.country_calling_code, offerCount: 0 });
    if (key === '') continue;
    const current = items.get(key);
    items.set(key, {
      countryISO2: current?.countryISO2 || offer.country_iso2,
      countryName: current?.countryName || offer.country_name,
      countryCallingCode: current?.countryCallingCode || offer.country_calling_code,
      offerCount: (current?.offerCount || 0) + 1
    });
  }
  return [...items.values()].sort((left, right) => countryChoiceLabel(left).localeCompare(countryChoiceLabel(right)));
}

export function filterAndSortOffers(offers: SmsPriceOffer[], providerKeys: string[], searchText: string, applicationKey: string, countryISO2: string, countryCallingCode: string, minAvailable: number, sort: OfferSort) {
  const providerFilter = new Set(providerKeys);
  const tokens = applicationKey ? [] : searchTokens(searchText);
  return offers
    .filter((offer) => providerFilter.size === 0 || providerFilter.has(offer.provider_key))
    .filter((offer) => applicationMatches(offer, applicationKey, searchText))
    .filter((offer) => countryMatches(offer, countryISO2, countryCallingCode))
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

export function applicationChoiceLabel(choice: ApplicationChoice) {
  return choice.displayName && choice.displayName !== choice.applicationKey ? `${choice.displayName} (${choice.applicationKey})` : choice.applicationKey;
}

export function matchApplicationChoice(value: string, choices: ApplicationChoice[]) {
  const normalized = value.trim().toLowerCase();
  return choices.find((choice) => [choice.applicationKey, choice.displayName, applicationChoiceLabel(choice)].some((item) => item.toLowerCase() === normalized));
}

export function countryChoiceValue(choice: CountryChoice) {
  if (!choice.countryISO2 && !choice.countryCallingCode) return '';
  return `${choice.countryISO2}|${choice.countryCallingCode}`;
}

export function countryChoiceLabel(choice: CountryChoice) {
  return [choice.countryName, choice.countryISO2, choice.countryCallingCode && `+${choice.countryCallingCode}`].filter(Boolean).join(' · ');
}

export function countryValue(iso2: string, callingCode: string) {
  return countryChoiceValue({ countryISO2: iso2, countryName: '', countryCallingCode: callingCode, offerCount: 0 });
}

export function parseCountryValue(value: string) {
  const [countryISO2 = '', countryCallingCode = ''] = value.split('|');
  return { countryISO2, countryCallingCode };
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

function providerOffers(offers: SmsPriceOffer[], providerKeys: string[]) {
  const providerFilter = new Set(providerKeys);
  return offers.filter((offer) => providerFilter.size === 0 || providerFilter.has(offer.provider_key));
}

function applicationMatches(offer: SmsPriceOffer, applicationKey: string, searchText: string) {
  if (applicationKey) return offer.application_key.toLowerCase() === applicationKey.toLowerCase();
  const tokens = searchTokens(searchText);
  if (tokens.length === 0) return true;
  const haystack = [offer.application_key, offer.application_name].filter(Boolean).join(' ').toLowerCase();
  return tokens.every((token) => haystack.includes(token));
}

function countryMatches(offer: SmsPriceOffer, countryISO2: string, countryCallingCode: string) {
  if (countryISO2 && offer.country_iso2.toLowerCase() !== countryISO2.toLowerCase()) return false;
  return !(countryCallingCode && offer.country_calling_code !== countryCallingCode);
}
