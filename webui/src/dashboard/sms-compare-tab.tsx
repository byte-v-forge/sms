import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { listSmsApplications, listSmsCountries, listSmsPriceOffers, smsKeys } from './sms-api';
import { SmsCompareForm } from './sms-compare-form';
import { bestOffer, enabledProviderKeys, filterAndSortOffers, providerChoices, type OfferSort } from './sms-compare-data';
import { applicationChoices, countryChoices, countryValue, matchApplicationChoice, parseCountryValue } from './sms-compare-options';
import { canSearch, compareQueryFromSearch, compareQuerySearchParams, smsPriceOfferQuery, type CompareQuery } from './sms-compare-query';
import { CompareSummary, OffersTable } from './sms-compare-table';

type SmsCompareTabProps = {
  providerOptions: SmsProviderPluginDescriptor[];
  configs: SmsProviderConfig[];
  acquiringOfferId?: string;
  onAcquire: (offer: SmsPriceOffer) => void;
};

export function SmsCompareTab({ providerOptions, configs, acquiringOfferId, onAcquire }: SmsCompareTabProps) {
  const [searchParams, setSearchParams] = useSearchParams();
  const choices = useMemo(() => providerChoices(providerOptions, configs), [configs, providerOptions]);
  const enabledKeys = useMemo(() => enabledProviderKeys(choices), [choices]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>();
  const [searchText, setSearchText] = useState('');
  const [applicationKey, setApplicationKey] = useState('');
  const [countryISO2, setCountryISO2] = useState('');
  const [countryCallingCode, setCountryCallingCode] = useState('');
  const [minAvailable, setMinAvailable] = useState(1);
  const [sort, setSort] = useState<OfferSort>('price');
  const searchKey = searchParams.toString();
  const routeQuery = useMemo(() => compareQueryFromSearch(new URLSearchParams(searchKey)), [searchKey]);

  useEffect(() => {
    setSearchText(routeQuery.searchText);
    setApplicationKey(routeQuery.applicationKey);
    setCountryISO2(routeQuery.countryISO2);
    setCountryCallingCode(routeQuery.countryCallingCode);
    setMinAvailable(routeQuery.minAvailable);
    setSort(routeQuery.sort);
    setSelectedKeys(routeQuery.providerKeys.length > 0 ? routeQuery.providerKeys : undefined);
  }, [routeQuery]);

  useEffect(() => {
    if (routeQuery.providerKeys.length === 0 && enabledKeys.length > 0) setSelectedKeys(enabledKeys);
  }, [enabledKeys, routeQuery.providerKeys.length]);

  const activeKeys = (selectedKeys || enabledKeys).filter((key) => enabledKeys.includes(key));
  const currentQuery: CompareQuery = { searchText: searchText.trim(), applicationKey, countryISO2, countryCallingCode, providerKeys: activeKeys, minAvailable: Math.max(0, minAvailable), sort };
  const serverQuery = smsPriceOfferQuery(currentQuery);
  const queried = canSearch(serverQuery.providerKeys);
  const applicationsQuery = useQuery({ queryKey: smsKeys.applications({ providerKeys: activeKeys }), queryFn: () => listSmsApplications({ providerKeys: activeKeys }), enabled: queried });
  const countriesQuery = useQuery({ queryKey: smsKeys.countries({ providerKeys: activeKeys, applicationKey: applicationKey || undefined }), queryFn: () => listSmsCountries({ providerKeys: activeKeys, applicationKey: applicationKey || undefined }), enabled: queried });
  const offersQuery = useQuery({ queryKey: smsKeys.priceOffers(serverQuery), queryFn: () => listSmsPriceOffers(serverQuery), enabled: queried });
  const allOffers = offersQuery.data?.offers || [];
  const applicationOptions = applicationChoices(applicationsQuery.data?.applications || [], allOffers);
  const countryOptions = countryChoices(countriesQuery.data?.countries || [], allOffers);
  const offers = filterAndSortOffers(allOffers, serverQuery.providerKeys, currentQuery.searchText, currentQuery.applicationKey, currentQuery.countryISO2, currentQuery.countryCallingCode, currentQuery.minAvailable, currentQuery.sort);
  const top = bestOffer(offers);
  const error = offersQuery.data?.error?.message;

  useEffect(() => {
    if (searchText || !applicationKey || applicationOptions.length === 0) return;
    const selected = applicationOptions.find((item) => item.applicationKey === applicationKey);
    if (selected) setSearchText(selected.displayName || selected.applicationKey);
  }, [applicationKey, applicationOptions, searchText]);

  function submitQuery(event?: FormEvent) {
    event?.preventDefault();
    setSearchParams(compareQuerySearchParams(draftQuery(sort, activeKeys)));
  }

  function changeSort(next: OfferSort) {
    setSort(next);
    const draft = draftQuery(next, activeKeys);
    if (canSearch(draft.providerKeys)) setSearchParams(compareQuerySearchParams(draft));
  }

  function draftQuery(nextSort: OfferSort, providerKeys: string[]): CompareQuery {
    const resolvedApplication = applicationKey || matchApplicationChoice(searchText, applicationOptions)?.applicationKey || '';
    return { searchText: searchText.trim(), applicationKey: resolvedApplication, countryISO2, countryCallingCode, providerKeys: [...providerKeys], minAvailable: Math.max(0, minAvailable), sort: nextSort };
  }

  function resetFilters() {
    setSearchParams(new URLSearchParams());
    setSearchText('');
    setApplicationKey('');
    setCountryISO2('');
    setCountryCallingCode('');
    setMinAvailable(1);
    setSort('price');
    setSelectedKeys(enabledKeys);
  }

  function changeServiceSearch(value: string) {
    setSearchText(value);
    if (applicationKey && value !== applicationLabel(applicationKey)) {
      setApplicationKey('');
      changeCountry('');
    }
  }

  function changeApplication(value: string) {
    const selected = applicationOptions.find((item) => item.applicationKey === value);
    setApplicationKey(value);
    setSearchText(selected ? selected.displayName || selected.applicationKey : '');
    if (value !== applicationKey) changeCountry('');
  }

  function changeCountry(value: string) {
    const parsed = parseCountryValue(value);
    setCountryISO2(parsed.countryISO2);
    setCountryCallingCode(parsed.countryCallingCode);
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col bg-muted/20">
      <SmsCompareForm
        choices={choices}
        applications={applicationOptions}
        countries={countryOptions}
        searchText={searchText}
        applicationKey={applicationKey}
        countryValue={countryValue(countryISO2, countryCallingCode)}
        providerKeys={activeKeys}
        minAvailable={minAvailable}
        sort={sort}
        canSubmit={canSearch(activeKeys)}
        onSearchTextChange={changeServiceSearch}
        onApplicationChange={changeApplication}
        onCountryChange={changeCountry}
        onProviderKeysChange={setSelectedKeys}
        onMinAvailableChange={setMinAvailable}
        onSortChange={changeSort}
        onSubmit={submitQuery}
        onReset={resetFilters}
      />
      <CompareSummary loading={offersQuery.isLoading} total={offers.length} providerCount={new Set(offers.map((offer) => offer.provider_key)).size} best={top} error={error} providerErrors={offersQuery.data?.provider_errors || []} />
      <OffersTable offers={offers} top={top} loading={offersQuery.isLoading} queried={queried} error={error} acquiringOfferId={acquiringOfferId} onAcquire={onAcquire} />
    </div>
  );

  function applicationLabel(key: string) {
    const selected = applicationOptions.find((item) => item.applicationKey === key);
    return selected?.displayName || selected?.applicationKey || '';
  }
}
