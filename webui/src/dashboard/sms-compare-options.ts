import type { SmsApplicationInfo, SmsCountry, SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';

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

export function applicationChoices(applications: SmsApplicationInfo[], offers: SmsPriceOffer[]) {
  const items = new Map<string, ApplicationChoice>();
  for (const app of applications) {
    upsertApplication(items, app.application_key, app.display_name, 0);
  }
  for (const offer of offers) {
    upsertApplication(items, offer.application_key, offer.application_name, 1);
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
  const normalized = normalizeChoiceText(value);
  return choices.find((choice) => [choice.applicationKey, choice.displayName, applicationChoiceLabel(choice)].some((item) => normalizeChoiceText(item) === normalized));
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

function upsertApplication(items: Map<string, ApplicationChoice>, key: string, name: string, offerCount: number) {
  key = key.trim();
  if (!key) return;
  const current = items.get(key);
  items.set(key, {
    applicationKey: key,
    displayName: bestDisplayName(current?.displayName, name, key),
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

function normalizeChoiceText(value: string) {
  return value.trim().toLowerCase();
}
