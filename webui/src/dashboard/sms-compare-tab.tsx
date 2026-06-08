import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { listSmsPriceOffers, smsKeys } from './sms-api';
import { SmsCompareForm } from './sms-compare-form';
import { bestOffer, enabledProviderKeys, filterAndSortOffers, providerChoices, type OfferSort } from './sms-compare-data';
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
  const [applicationKey, setApplicationKey] = useState('');
  const [countryISO2, setCountryISO2] = useState('');
  const [countryCallingCode, setCountryCallingCode] = useState('');
  const [minAvailable, setMinAvailable] = useState(1);
  const [sort, setSort] = useState<OfferSort>('price');
  const searchKey = searchParams.toString();
  const routeQuery = useMemo(() => compareQueryFromSearch(new URLSearchParams(searchKey)), [searchKey]);

  useEffect(() => {
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
  const compareQuery = useMemo(() => activeCompareQuery(routeQuery, enabledKeys), [enabledKeys, routeQuery]);
  const serverQuery = useMemo(() => smsPriceOfferQuery(compareQuery), [compareQuery]);
  const queried = canSearch(serverQuery.applicationKey, serverQuery.countryISO2, serverQuery.countryCallingCode, serverQuery.providerKeys);
  const offersQuery = useQuery({ queryKey: smsKeys.priceOffers(serverQuery), queryFn: () => listSmsPriceOffers(serverQuery), enabled: queried });
  const offers = filterAndSortOffers(offersQuery.data?.offers || [], serverQuery.providerKeys, serverQuery.minAvailable, compareQuery.sort);
  const top = bestOffer(offers);
  const error = offersQuery.data?.error?.message;

  function submitQuery(event?: FormEvent) {
    event?.preventDefault();
    setSearchParams(compareQuerySearchParams(draftQuery(sort, activeKeys)));
  }

  function changeSort(next: OfferSort) {
    setSort(next);
    const draft = draftQuery(next, activeKeys);
    if (canSearch(draft.applicationKey, draft.countryISO2, draft.countryCallingCode, draft.providerKeys)) setSearchParams(compareQuerySearchParams(draft));
  }

  function draftQuery(nextSort: OfferSort, providerKeys: string[]): CompareQuery {
    return { applicationKey: applicationKey.trim(), countryISO2: countryISO2.trim().toUpperCase(), countryCallingCode: countryCallingCode.trim().replace(/^\+/, ''), providerKeys: [...providerKeys], minAvailable: Math.max(0, minAvailable), sort: nextSort };
  }

  function resetFilters() {
    setSearchParams(new URLSearchParams());
    setApplicationKey('');
    setCountryISO2('');
    setCountryCallingCode('');
    setMinAvailable(1);
    setSort('price');
    setSelectedKeys(enabledKeys);
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col bg-muted/20">
      <SmsCompareForm
        choices={choices}
        applicationKey={applicationKey}
        countryISO2={countryISO2}
        countryCallingCode={countryCallingCode}
        providerKeys={activeKeys}
        minAvailable={minAvailable}
        sort={sort}
        canSubmit={canSearch(applicationKey, countryISO2, countryCallingCode, activeKeys)}
        onApplicationKeyChange={setApplicationKey}
        onCountryISO2Change={setCountryISO2}
        onCountryCallingCodeChange={setCountryCallingCode}
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

function activeCompareQuery(routeQuery: CompareQuery, enabledKeys: string[]): CompareQuery {
  return {
    ...routeQuery,
    providerKeys: routeQuery.providerKeys.length > 0 ? routeQuery.providerKeys.filter((key) => enabledKeys.includes(key)) : enabledKeys
  };
}
