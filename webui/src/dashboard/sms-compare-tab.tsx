import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { Search } from 'lucide-react';
import { useSearchParams } from 'react-router';
import { Badge, Button, Input, useQuery } from '../ui';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { listSmsPriceOffers, smsKeys, type SmsPriceOfferQuery, type SmsProviderOption, type SmsProviderSetting } from './sms-api';
import {
  bestOffer,
  enabledProviderKeys,
  filterAndSortOffers,
  providerChoices,
  type OfferSort,
  type ProviderChoice
} from './sms-compare-data';
import { CompareSummary, OffersTable } from './sms-compare-table';

type CompareQuery = {
  applicationKey: string;
  countryISO2: string;
  countryCallingCode: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
};

type SmsCompareTabProps = {
  providerOptions: SmsProviderOption[];
  configs: SmsProviderSetting[];
  acquiringOfferId?: string;
  onAcquire: (offer: SmsPriceOffer) => void;
};

type ProviderPickerProps = {
  choices: ProviderChoice[];
  selectedKeys: string[];
  onChange: (keys: string[]) => void;
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
    if (!selectedKeys && enabledKeys.length > 0) {
      setSelectedKeys(enabledKeys);
    }
  }, [enabledKeys, selectedKeys]);
  const activeKeys = selectedKeys || enabledKeys;
  const compareQuery = useMemo(() => ({ ...routeQuery, providerKeys: routeQuery.providerKeys.length > 0 ? routeQuery.providerKeys : enabledKeys }), [enabledKeys, routeQuery]);
  const serverQuery = useMemo(() => smsPriceOfferQuery(compareQuery), [compareQuery]);
  const queried = canSearch(serverQuery.applicationKey, serverQuery.countryISO2, serverQuery.countryCallingCode, serverQuery.providerKeys);
  const offersQuery = useQuery({ queryKey: smsKeys.priceOffers(serverQuery), queryFn: () => listSmsPriceOffers(serverQuery), enabled: queried });
  const offers = filterAndSortOffers(offersQuery.data?.offers || [], serverQuery.providerKeys, serverQuery.minAvailable, compareQuery.sort);
  const top = bestOffer(offers);
  const error = offersQuery.data?.error?.message;

  function submitQuery(event?: FormEvent) {
    event?.preventDefault();
    setSearchParams(compareQuerySearchParams({
      applicationKey: applicationKey.trim(),
      countryISO2: countryISO2.trim().toUpperCase(),
      countryCallingCode: countryCallingCode.trim().replace(/^\+/, ''),
      providerKeys: [...activeKeys],
      minAvailable: Math.max(0, minAvailable),
      sort
    }));
  }

  function changeSort(next: OfferSort) {
    setSort(next);
    if (queried) setSearchParams(compareQuerySearchParams({ ...compareQuery, sort: next }));
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col bg-muted/20">
      <form className="grid gap-3 border-b border-border bg-card p-3 lg:grid-cols-[1fr_auto]" onSubmit={submitQuery}>
        <div className="grid gap-2 md:grid-cols-4">
          <Input aria-label="应用" placeholder="应用，如 whatsapp/gojek" value={applicationKey} onChange={(event) => setApplicationKey(event.target.value)} />
          <Input aria-label="国家 ISO2" placeholder="国家 ISO2，如 ID" value={countryISO2} onChange={(event) => setCountryISO2(event.target.value)} />
          <Input aria-label="国家区号" placeholder="区号，如 62" value={countryCallingCode} onChange={(event) => setCountryCallingCode(event.target.value)} />
          <Input aria-label="最低库存" min={0} type="number" value={minAvailable} onChange={(event) => setMinAvailable(numberInputValue(event.target.value))} />
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <select className="h-8 rounded-lg border border-border bg-background px-2 text-xs" value={sort} onChange={(event) => changeSort(event.target.value as OfferSort)}>
            <option value="price">按低价</option>
            <option value="available">按库存</option>
            <option value="provider">按平台</option>
          </select>
          <Button disabled={!canSearch(applicationKey, countryISO2, countryCallingCode, activeKeys)} type="submit">
            <Search className="mr-1 size-4" />查询比对
          </Button>
        </div>
        <ProviderPicker choices={choices} selectedKeys={activeKeys} onChange={setSelectedKeys} />
      </form>
      <CompareSummary loading={offersQuery.isLoading} total={offers.length} providerCount={new Set(offers.map((offer) => offer.provider_key)).size} best={top} error={error} />
      <OffersTable offers={offers} top={top} loading={offersQuery.isLoading} queried={queried} error={error} acquiringOfferId={acquiringOfferId} onAcquire={onAcquire} />
    </div>
  );
}

function ProviderPicker({ choices, selectedKeys, onChange }: ProviderPickerProps) {
  return (
    <div className="flex flex-wrap gap-2 lg:col-span-2">
      {choices.map((choice) => {
        const checked = selectedKeys.includes(choice.providerKey);
        return (
          <label key={choice.providerKey} className={`inline-flex h-8 items-center gap-2 rounded-lg border px-2 text-xs ${choice.enabled ? 'bg-background' : 'opacity-50'}`}>
            <input checked={checked} disabled={!choice.enabled} type="checkbox" onChange={() => onChange(checked ? selectedKeys.filter((key) => key !== choice.providerKey) : [...selectedKeys, choice.providerKey])} />
            <span>{choice.displayName}</span>
            {!choice.configured && <Badge variant="outline">未配置</Badge>}
          </label>
        );
      })}
      {choices.length === 0 && <span className="text-xs text-muted-foreground">暂无 provider 插件</span>}
    </div>
  );
}

function canSearch(applicationKey: string, countryISO2: string, callingCode: string, providerKeys: string[]) {
  return applicationKey.trim() !== '' && (countryISO2.trim() !== '' || callingCode.trim() !== '') && providerKeys.length > 0;
}

function numberInputValue(value: string) {
  const number = Number(value || 0);
  return Number.isFinite(number) ? number : 0;
}

function compareQueryFromSearch(params: URLSearchParams): CompareQuery {
  return {
    applicationKey: params.get('application_key') || '',
    countryISO2: params.get('country_iso2') || '',
    countryCallingCode: params.get('country_calling_code') || '',
    providerKeys: params.getAll('provider_key').filter(Boolean),
    minAvailable: Math.max(0, numberInputValue(params.get('min_available') || '1')),
    sort: offerSort(params.get('sort'))
  };
}

function smsPriceOfferQuery(query: CompareQuery): SmsPriceOfferQuery {
  return {
    applicationKey: query.applicationKey,
    countryISO2: query.countryISO2,
    countryCallingCode: query.countryCallingCode,
    providerKeys: query.providerKeys,
    minAvailable: query.minAvailable
  };
}

function compareQuerySearchParams(query: CompareQuery) {
  const params = new URLSearchParams();
  params.set('application_key', query.applicationKey);
  if (query.countryISO2) params.set('country_iso2', query.countryISO2);
  if (query.countryCallingCode) params.set('country_calling_code', query.countryCallingCode);
  for (const providerKey of query.providerKeys) params.append('provider_key', providerKey);
  params.set('min_available', String(query.minAvailable));
  params.set('sort', query.sort);
  return params;
}

function offerSort(value: string | null): OfferSort {
  if (value === 'available' || value === 'provider') return value;
  return 'price';
}
