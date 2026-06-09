import type { SmsCountry, SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderApplications } from './sms-api';
import { normalizeChoiceToken, normalizeDisplayName, searchTerms, termsContain, titleSlug } from './sms-compare-text';
import type { SearchSelectOption } from './sms-search-select';

export type ApplicationChoice = {
  id: string;
  displayName: string;
  aliases: string[];
  offerCount: number;
  providerApplicationKeys: Record<string, string>;
};

export type CountryChoice = {
  countryISO2: string;
  countryName: string;
  countryCallingCode: string;
  offerCount: number;
};

export function applicationChoices(results: SmsProviderApplications[], offers: SmsPriceOffer[]) {
  const items = new Map<string, ApplicationChoice>();
  for (const result of results) {
    for (const app of result.applications) {
      upsertApplication(items, result.providerKey, app.application_key, app.display_name, 0, app.aliases || []);
    }
  }
  for (const offer of offers) {
    const routeKey = offer.offer_ref?.route_ref?.upstream_service_key || offer.application_key;
    upsertApplication(items, offer.provider_key, routeKey, offer.application_name, 1, [offer.application_key]);
  }
  return [...items.values()].sort((left, right) => applicationChoiceLabel(left).localeCompare(applicationChoiceLabel(right)));
}

export function countryChoices(countries: SmsCountry[], offers: SmsPriceOffer[]) {
  const items = new Map<string, CountryChoice>();
  for (const country of countries) {
    upsertCountry(items, country.country_iso2, country.name, country.country_calling_code, 0);
  }
  for (const offer of offers) {
    upsertCountry(items, offer.country_iso2, offer.country_name, offer.country_calling_code, 1);
  }
  return [...items.values()].sort((left, right) => countryChoiceLabel(left).localeCompare(countryChoiceLabel(right)));
}

export function applicationChoiceLabel(choice: ApplicationChoice) {
  return choice.displayName || '未命名服务';
}

export function matchApplicationChoice(value: string, choices: ApplicationChoice[]) {
  const normalized = normalizeChoiceToken(value);
  if (!normalized) return undefined;
  const exact = choices.find((choice) => applicationSearchTokens(choice).some((item) => item === normalized));
  if (exact) return exact;
  const primaryExact = choices.filter((choice) => normalizeChoiceToken(primaryApplicationName(choice.displayName)) === normalized);
  if (primaryExact.length === 1) return primaryExact[0];
  const partial = choices.filter((choice) => applicationChoiceMatches(choice, value));
  return partial.length === 1 ? partial[0] : undefined;
}

export function providerApplicationKey(choice: ApplicationChoice | undefined, providerKey: string, fallback: string) {
  return choice?.providerApplicationKeys[providerKey] || choice?.displayName || fallback.trim();
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

export function applicationSelectOptions(choices: ApplicationChoice[]): SearchSelectOption[] {
  return choices.map((choice) => ({
    value: choice.id,
    label: applicationChoiceLabel(choice),
    description: providerCountDescription(choice),
    badge: choice.offerCount > 0 ? String(choice.offerCount) : undefined,
    keywords: applicationSearchTokens(choice)
  }));
}

export function countrySelectOptions(choices: CountryChoice[]): SearchSelectOption[] {
  return choices.map((choice) => ({
    value: countryChoiceValue(choice),
    label: choice.countryName || choice.countryISO2 || `+${choice.countryCallingCode}`,
    description: [choice.countryISO2, choice.countryCallingCode && `+${choice.countryCallingCode}`].filter(Boolean).join(' · '),
    badge: choice.offerCount > 0 ? String(choice.offerCount) : undefined,
    keywords: [choice.countryName, choice.countryISO2, choice.countryCallingCode, `+${choice.countryCallingCode}`].filter(Boolean)
  }));
}

function upsertApplication(items: Map<string, ApplicationChoice>, providerKey: string, key: string, name: string, offerCount: number, aliases: string[]) {
  const displayName = applicationDisplayName(key, name);
  if (!displayName) return;
  const identity = applicationIdentity(displayName, key);
  const current = items.get(identity);
  items.set(identity, {
    id: identity,
    displayName: bestDisplayName(current?.displayName, displayName),
    aliases: uniqueTexts([...(current?.aliases || []), key, displayName, ...aliases]),
    offerCount: (current?.offerCount || 0) + offerCount,
    providerApplicationKeys: providerApplicationKeys(current, providerKey, key)
  });
}

function upsertCountry(items: Map<string, CountryChoice>, iso2: string, name: string, callingCode: string, offerCount: number) {
  iso2 = iso2.trim().toUpperCase();
  callingCode = callingCode.trim().replace(/^\+/, '');
  const item = { countryISO2: iso2, countryName: name, countryCallingCode: callingCode, offerCount: 0 };
  const key = countryChoiceValue(item);
  if (!key) return;
  const current = items.get(key);
  items.set(key, {
    countryISO2: current?.countryISO2 || iso2,
    countryName: bestDisplayName(current?.countryName, name || iso2 || callingCode),
    countryCallingCode: current?.countryCallingCode || callingCode,
    offerCount: (current?.offerCount || 0) + offerCount
  });
}

function applicationDisplayName(key: string, name: string) {
  const display = normalizeDisplayName(name);
  if (display && !looksLikeShortCode(display, key)) return display;
  const fromKey = titleSlug(key);
  return fromKey.length > 3 ? fromKey : '';
}

function applicationIdentity(displayName: string, key: string) {
  return normalizeChoiceToken(primaryApplicationName(displayName)) || normalizeChoiceToken(displayName) || normalizeChoiceToken(key);
}

function primaryApplicationName(value: string) {
  return value.split(/[,/|;]+/)[0]?.trim() || value.trim();
}

function bestDisplayName(current = '', candidate = '') {
  const name = normalizeDisplayName(candidate);
  if (!current || (name && name.length > current.length)) return name;
  return current;
}

function providerApplicationKeys(current: ApplicationChoice | undefined, providerKey: string, key: string) {
  const values = { ...(current?.providerApplicationKeys || {}) };
  if (providerKey && key && !values[providerKey]) values[providerKey] = key.trim();
  return values;
}

function applicationChoiceMatches(choice: ApplicationChoice, value: string) {
  const queryTerms = searchTerms(value);
  if (queryTerms.length === 0) return false;
  return [choice.displayName, ...choice.aliases].some((candidate) => termsContain(searchTerms(candidate), queryTerms));
}

function applicationSearchTokens(choice: ApplicationChoice) {
  return [choice.displayName, ...choice.aliases].map(normalizeChoiceToken).filter(Boolean);
}

function providerCountDescription(choice: ApplicationChoice) {
  const count = Object.keys(choice.providerApplicationKeys).length;
  return count > 1 ? `${count} 个平台` : undefined;
}

function looksLikeShortCode(displayName: string, key: string) {
  return displayName.length <= 3 && displayName === displayName.toLowerCase() && normalizeChoiceToken(displayName) === normalizeChoiceToken(key);
}

function uniqueTexts(values: string[]) {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const value of values) {
    const text = normalizeDisplayName(value);
    const token = normalizeChoiceToken(text);
    if (!text || !token || seen.has(token)) continue;
    seen.add(token);
    out.push(text);
  }
  return out;
}
