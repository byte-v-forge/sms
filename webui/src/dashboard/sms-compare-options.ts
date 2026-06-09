import type { SmsApplicationInfo, SmsCountry, SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SearchSelectOption } from './sms-search-select';

export type ApplicationChoice = {
  applicationKey: string;
  displayName: string;
  aliases: string[];
  offerCount: number;
};

export type CountryChoice = {
  countryISO2: string;
  countryName: string;
  countryCallingCode: string;
  offerCount: number;
};

export function applicationChoices(applications: SmsApplicationInfo[], offers: SmsPriceOffer[]) {
  const items = new Map<string, ApplicationChoice>();
  for (const app of applications) {
    upsertApplication(items, app.application_key, app.display_name, 0, app.aliases || []);
  }
  for (const offer of offers) {
    upsertApplication(items, offer.application_key, offer.application_name, 1, [offer.offer_ref?.route_ref?.upstream_service_key || '']);
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
  return choice.displayName || choice.applicationKey;
}

export function matchApplicationChoice(value: string, choices: ApplicationChoice[]) {
  const normalized = normalizeChoiceToken(value);
  if (!normalized) return undefined;
  const exact = choices.find((choice) => applicationSearchTokens(choice).some((item) => item === normalized));
  if (exact) return exact;
  const partial = choices.filter((choice) => applicationSearchTokens(choice).some((item) => item.includes(normalized)));
  return partial.length === 1 ? partial[0] : undefined;
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
    value: choice.applicationKey,
    label: applicationChoiceLabel(choice),
    description: applicationDescription(choice),
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

function upsertApplication(items: Map<string, ApplicationChoice>, key: string, name: string, offerCount: number, aliases: string[]) {
  key = key.trim();
  if (!key) return;
  const displayName = name.trim() || key;
  const identity = normalizeChoiceToken(displayName) || normalizeChoiceToken(key);
  const current = items.get(identity);
  items.set(identity, {
    applicationKey: bestApplicationKey(identity, current?.applicationKey, key, displayName),
    displayName: bestDisplayName(current?.displayName, displayName, key),
    aliases: uniqueTexts([...(current?.aliases || []), key, name, ...aliases]),
    offerCount: (current?.offerCount || 0) + offerCount
  });
}

function upsertCountry(items: Map<string, CountryChoice>, iso2: string, name: string, callingCode: string, offerCount: number) {
  const item = { countryISO2: iso2, countryName: name, countryCallingCode: callingCode, offerCount: 0 };
  const key = countryChoiceValue(item);
  if (!key) return;
  const current = items.get(key);
  items.set(key, {
    countryISO2: current?.countryISO2 || iso2,
    countryName: bestDisplayName(current?.countryName, name, iso2 || callingCode),
    countryCallingCode: current?.countryCallingCode || callingCode,
    offerCount: (current?.offerCount || 0) + offerCount
  });
}

function bestDisplayName(current = '', candidate = '', fallback = '') {
  const name = candidate.trim() || fallback.trim();
  if (!current || (name && name.length > current.length)) return name;
  return current;
}

function applicationSearchTokens(choice: ApplicationChoice) {
  return [choice.applicationKey, choice.displayName, applicationChoiceLabel(choice), ...choice.aliases].map(normalizeChoiceToken).filter(Boolean);
}

function normalizeChoiceToken(value: string) {
  return value.trim().toLowerCase().replace(/[^a-z0-9]+/g, '');
}

function applicationDescription(choice: ApplicationChoice) {
  const baseTokens = new Set([choice.applicationKey, choice.displayName].map(normalizeChoiceToken));
  const aliases = choice.aliases.filter((alias) => !baseTokens.has(normalizeChoiceToken(alias))).slice(0, 3);
  return aliases.length > 0 ? aliases.join(' / ') : undefined;
}

function bestApplicationKey(identity: string, current = '', candidate = '', displayName = '') {
  if (candidate.trim() === identity) return candidate;
  if (current.trim() === identity) return current;
  if (normalizeChoiceToken(current) === identity) return current;
  if (normalizeChoiceToken(candidate) === identity) return candidate;
  return displayName || current || candidate;
}

function uniqueTexts(values: string[]) {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const value of values) {
    const text = value.trim();
    const token = normalizeChoiceToken(text);
    if (!text || !token || seen.has(token)) continue;
    seen.add(token);
    out.push(text);
  }
  return out;
}
