import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { listSmsApplicationsByProvider, listSmsCountries, listSmsPriceOffersByProvider, smsKeys } from './sms-api';
import { SmsCompareForm } from './sms-compare-form';
import { bestOffer, enabledProviderKeys, filterAndSortOffers, providerChoices, type OfferSort } from './sms-compare-data';
import { providerOfferQueries } from './sms-compare-offer-queries';
import { applicationChoices, countryChoices, countryValue, matchApplicationChoice, parseCountryValue, type ApplicationChoice } from './sms-compare-options';
import { canSearch, canUseCatalog, compareQueryFromSearch, compareQuerySearchParams, type CompareQuery } from './sms-compare-query';
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
  const [serviceText, setServiceText] = useState('');
  const [applicationId, setApplicationId] = useState('');
  const [applicationSnapshot, setApplicationSnapshot] = useState<ApplicationChoice>();
  const [countryISO2, setCountryISO2] = useState('');
  const [countryCallingCode, setCountryCallingCode] = useState('');
  const [minAvailable, setMinAvailable] = useState(1);
  const [sort, setSort] = useState<OfferSort>('price');
  const searchKey = searchParams.toString();
  const routeQuery = useMemo(() => compareQueryFromSearch(new URLSearchParams(searchKey)), [searchKey]);

  useEffect(() => {
    setServiceText(routeQuery.serviceText);
    setApplicationId('');
    setApplicationSnapshot(undefined);
    setCountryISO2(routeQuery.countryISO2);
    setCountryCallingCode(routeQuery.countryCallingCode);
    setMinAvailable(routeQuery.minAvailable);
    setSort(routeQuery.sort);
  }, [routeQuery]);

  useEffect(() => {
    if (routeQuery.providerKeys.length === 0 && enabledKeys.length > 0) setSelectedKeys(enabledKeys);
  }, [enabledKeys, routeQuery.providerKeys.length]);

  const activeKeys = (selectedKeys || enabledKeys).filter((key) => enabledKeys.includes(key));
  const currentQuery: CompareQuery = { serviceText: serviceText.trim(), countryISO2, countryCallingCode, providerKeys: activeKeys, minAvailable: Math.max(0, minAvailable), sort };
  const catalogEnabled = canUseCatalog(activeKeys);
  const applicationsQuery = useQuery({ queryKey: smsKeys.providerApplications({ providerKeys: activeKeys, searchText: currentQuery.serviceText }), queryFn: () => listSmsApplicationsByProvider({ providerKeys: activeKeys, searchText: currentQuery.serviceText }), enabled: catalogEnabled && currentQuery.serviceText.length > 0 });
  const countriesQuery = useQuery({ queryKey: smsKeys.countries({ providerKeys: activeKeys }), queryFn: () => listSmsCountries({ providerKeys: activeKeys }), enabled: catalogEnabled });
  const applicationOptions = applicationChoices(applicationsQuery.data || [], []);
  const optionApplication = applicationOptions.find((item) => item.id === applicationId) || matchApplicationChoice(currentQuery.serviceText, applicationOptions);
  const selectedApplication = selectedApplicationChoice(optionApplication, applicationSnapshot, applicationId);
  const offerRequests = providerOfferQueries(currentQuery, selectedApplication);
  const applicationCatalogReady = !applicationsQuery.isLoading && !applicationsQuery.isFetching;
  const offersQuery = useQuery({ queryKey: smsKeys.providerPriceOffers(offerRequests), queryFn: () => listSmsPriceOffersByProvider(offerRequests), enabled: applicationCatalogReady && canSearch(activeKeys, currentQuery.serviceText) && offerRequests.length > 0 });
  const allOffers = offersQuery.data?.offers || [];
  const mergedApplicationOptions = applicationChoices(applicationsQuery.data || [], allOffers);
  const countryOptions = countryChoices(countriesQuery.data?.countries || [], allOffers);
  const offers = filterAndSortOffers(allOffers, activeKeys, currentQuery.countryISO2, currentQuery.countryCallingCode, currentQuery.minAvailable, currentQuery.sort);
  const top = bestOffer(offers);
  const error = offersQuery.data?.error?.message;

  useEffect(() => {
    if (applicationId || !currentQuery.serviceText || mergedApplicationOptions.length === 0) return;
    const selected = matchApplicationChoice(currentQuery.serviceText, mergedApplicationOptions);
    if (selected) {
      setApplicationId(selected.id);
      setApplicationSnapshot(selected);
      setServiceText(selected.displayName);
    }
  }, [applicationId, currentQuery.serviceText, mergedApplicationOptions]);

  function submitQuery(event?: FormEvent) {
    event?.preventDefault();
    setSearchParams(compareQuerySearchParams(draftQuery(sort, activeKeys)));
  }

  function changeSort(next: OfferSort) {
    setSort(next);
    const draft = draftQuery(next, activeKeys);
    if (canSearch(draft.providerKeys, draft.serviceText)) setSearchParams(compareQuerySearchParams(draft));
  }

  function draftQuery(nextSort: OfferSort, providerKeys: string[]): CompareQuery {
    const selected = mergedApplicationOptions.find((item) => item.id === applicationId) || matchApplicationChoice(serviceText, mergedApplicationOptions);
    return { serviceText: (selected?.displayName || serviceText).trim(), countryISO2, countryCallingCode, providerKeys: [...providerKeys], minAvailable: Math.max(0, minAvailable), sort: nextSort };
  }

  function resetFilters() {
    setSearchParams(new URLSearchParams());
    setServiceText('');
    setApplicationId('');
    setApplicationSnapshot(undefined);
    setCountryISO2('');
    setCountryCallingCode('');
    setMinAvailable(1);
    setSort('price');
    setSelectedKeys(enabledKeys);
  }

  function changeServiceSearch(value: string) {
    setServiceText(value);
    setApplicationId('');
    setApplicationSnapshot(undefined);
  }

  function changeApplication(value: string) {
    const selected = mergedApplicationOptions.find((item) => item.id === value);
    setApplicationId(value);
    setApplicationSnapshot(selected);
    setServiceText(selected?.displayName || '');
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
        applications={mergedApplicationOptions}
        countries={countryOptions}
        serviceText={serviceText}
        applicationId={applicationId}
        countryValue={countryValue(countryISO2, countryCallingCode)}
        providerKeys={activeKeys}
        minAvailable={minAvailable}
        sort={sort}
        canSubmit={canSearch(activeKeys, serviceText)}
        onServiceTextChange={changeServiceSearch}
        onApplicationChange={changeApplication}
        onCountryChange={changeCountry}
        onProviderKeysChange={setSelectedKeys}
        onMinAvailableChange={setMinAvailable}
        onSortChange={changeSort}
        onSubmit={submitQuery}
        onReset={resetFilters}
      />
      <CompareSummary loading={offersQuery.isLoading} total={offers.length} providerCount={new Set(offers.map((offer) => offer.provider_key)).size} best={top} error={error} providerErrors={offersQuery.data?.provider_errors || []} />
      <OffersTable offers={offers} top={top} loading={offersQuery.isLoading} queried={catalogEnabled} error={error} acquiringOfferId={acquiringOfferId} serviceName={currentQuery.serviceText} onAcquire={onAcquire} />
    </div>
  );
}

function selectedApplicationChoice(option: ApplicationChoice | undefined, snapshot: ApplicationChoice | undefined, applicationId: string) {
  if (!snapshot || snapshot.id !== applicationId) return option;
  if (!option) return snapshot;
  return {
    ...option,
    aliases: [...new Set([...snapshot.aliases, ...option.aliases])],
    providerApplicationKeys: { ...snapshot.providerApplicationKeys, ...option.providerApplicationKeys }
  };
}
