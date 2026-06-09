import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { listSmsPriceOffers, smsKeys } from './sms-api';
import { SmsCompareForm } from './sms-compare-form';
import { applicationChoices, bestOffer, countryChoices, countryValue, enabledProviderKeys, filterAndSortOffers, matchApplicationChoice, parseCountryValue, providerChoices, type OfferSort } from './sms-compare-data';
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
  const offersQuery = useQuery({ queryKey: smsKeys.priceOffers(serverQuery), queryFn: () => listSmsPriceOffers(serverQuery), enabled: queried });
  const allOffers = offersQuery.data?.offers || [];
  const applicationOptions = applicationChoices(allOffers, activeKeys);
  const countryOptions = countryChoices(allOffers, activeKeys, applicationKey, searchText);
  const offers = filterAndSortOffers(allOffers, serverQuery.providerKeys, currentQuery.searchText, currentQuery.applicationKey, currentQuery.countryISO2, currentQuery.countryCallingCode, currentQuery.minAvailable, currentQuery.sort);
  const top = bestOffer(offers);
  const error = offersQuery.data?.error?.message;

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
    return { searchText: searchText.trim(), applicationKey, countryISO2, countryCallingCode, providerKeys: [...providerKeys], minAvailable: Math.max(0, minAvailable), sort: nextSort };
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

  function changeService(value: string) {
    setSearchText(value);
    setApplicationKey(matchApplicationChoice(value, applicationOptions)?.applicationKey || '');
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
        countryValue={countryValue(countryISO2, countryCallingCode)}
        providerKeys={activeKeys}
        minAvailable={minAvailable}
        sort={sort}
        canSubmit={canSearch(activeKeys)}
        onSearchTextChange={changeService}
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
}
