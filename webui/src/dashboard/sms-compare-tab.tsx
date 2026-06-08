import { useEffect, useMemo, useState } from 'react';
import type { FormEvent, ReactNode } from 'react';
import { CheckSquare, RotateCcw, Search, Square } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router';
import { Badge, Button, Card, Checkbox, Input, Select } from '../ui';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { listSmsPriceOffers, smsKeys, type SmsPriceOfferQuery } from './sms-api';
import { bestOffer, enabledProviderKeys, filterAndSortOffers, providerChoices, type OfferSort, type ProviderChoice } from './sms-compare-data';
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
  const compareQuery = useMemo(() => ({
    ...routeQuery,
    providerKeys: routeQuery.providerKeys.length > 0 ? routeQuery.providerKeys.filter((key) => enabledKeys.includes(key)) : enabledKeys
  }), [enabledKeys, routeQuery]);
  const serverQuery = useMemo(() => smsPriceOfferQuery(compareQuery), [compareQuery]);
  const queried = canSearch(serverQuery.applicationKey, serverQuery.countryISO2, serverQuery.countryCallingCode, serverQuery.providerKeys);
  const offersQuery = useQuery({ queryKey: smsKeys.priceOffers(serverQuery), queryFn: () => listSmsPriceOffers(serverQuery), enabled: queried });
  const offers = filterAndSortOffers(offersQuery.data?.offers || [], serverQuery.providerKeys, serverQuery.minAvailable, compareQuery.sort);
  const providerErrors = offersQuery.data?.provider_errors || [];
  const top = bestOffer(offers);
  const error = offersQuery.data?.error?.message;

  function submitQuery(event?: FormEvent) {
    event?.preventDefault();
    setSearchParams(compareQuerySearchParams(draftQuery(sort)));
  }

  function changeSort(next: OfferSort) {
    setSort(next);
    const draft = draftQuery(next);
    if (canSearch(draft.applicationKey, draft.countryISO2, draft.countryCallingCode, draft.providerKeys)) setSearchParams(compareQuerySearchParams(draft));
  }

  function draftQuery(nextSort: OfferSort): CompareQuery {
    return { applicationKey: applicationKey.trim(), countryISO2: countryISO2.trim().toUpperCase(), countryCallingCode: countryCallingCode.trim().replace(/^\+/, ''), providerKeys: [...activeKeys], minAvailable: Math.max(0, minAvailable), sort: nextSort };
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
      <Card className="m-4 mb-0 p-3">
        <form className="grid gap-3" onSubmit={submitQuery}>
          <div className="grid gap-2 md:grid-cols-6">
            <Field label="应用" className="md:col-span-2"><Input placeholder="whatsapp/gojek" value={applicationKey} onChange={(event) => setApplicationKey(event.target.value)} /></Field>
            <Field label="国家 ISO2"><Input placeholder="ID" value={countryISO2} onChange={(event) => setCountryISO2(event.target.value)} /></Field>
            <Field label="国家区号"><Input placeholder="62" value={countryCallingCode} onChange={(event) => setCountryCallingCode(event.target.value)} /></Field>
            <Field label="最低库存"><Input min={0} type="number" value={minAvailable} onChange={(event) => setMinAvailable(numberInputValue(event.target.value))} /></Field>
            <Field label="排序"><Select value={sort} onChange={(event) => changeSort(event.target.value as OfferSort)}><option value="price">按低价</option><option value="available">按库存</option><option value="provider">按平台</option></Select></Field>
          </div>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <ProviderPicker choices={choices} selectedKeys={activeKeys} onChange={setSelectedKeys} />
            <div className="flex items-center gap-2">
              <Button aria-label="重置查询条件" title="重置查询条件" size="icon-sm" variant="outline" onClick={resetFilters}><RotateCcw className="size-4" /></Button>
              <Button disabled={!canSearch(applicationKey, countryISO2, countryCallingCode, activeKeys)} type="submit"><Search className="size-4" />查询比对</Button>
            </div>
          </div>
          <SearchHint applicationKey={applicationKey} countryISO2={countryISO2} countryCallingCode={countryCallingCode} providerKeys={activeKeys} />
        </form>
      </Card>
      <CompareSummary loading={offersQuery.isLoading} total={offers.length} providerCount={new Set(offers.map((offer) => offer.provider_key)).size} best={top} error={error} providerErrors={providerErrors} />
      <OffersTable offers={offers} top={top} loading={offersQuery.isLoading} queried={queried} error={error} acquiringOfferId={acquiringOfferId} onAcquire={onAcquire} />
    </div>
  );
}

function Field({ label, className, children }: { label: string; className?: string; children: ReactNode }) {
  return <label className={`grid gap-1 text-xs font-medium text-muted-foreground ${className || ''}`}><span>{label}</span>{children}</label>;
}

function ProviderPicker({ choices, selectedKeys, onChange }: { choices: ProviderChoice[]; selectedKeys: string[]; onChange: (keys: string[]) => void }) {
  if (choices.length === 0) return <span className="text-xs text-muted-foreground">暂无 provider 插件</span>;
  const enabledKeys = choices.filter((choice) => choice.enabled).map((choice) => choice.providerKey);
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Badge variant="secondary">已选 {selectedKeys.length}/{enabledKeys.length}</Badge>
      <Button aria-label="选择全部可用平台" title="选择全部可用平台" size="icon-sm" variant="ghost" onClick={() => onChange(enabledKeys)}><CheckSquare className="size-4" /></Button>
      <Button aria-label="清空平台选择" title="清空平台选择" size="icon-sm" variant="ghost" onClick={() => onChange([])}><Square className="size-4" /></Button>
      {choices.map((choice) => <ProviderChoiceItem key={choice.providerKey} choice={choice} selectedKeys={selectedKeys} onChange={onChange} />)}
    </div>
  );
}

function ProviderChoiceItem({ choice, selectedKeys, onChange }: { choice: ProviderChoice; selectedKeys: string[]; onChange: (keys: string[]) => void }) {
  const checked = selectedKeys.includes(choice.providerKey);
  return (
    <label className={`inline-flex h-8 items-center gap-2 rounded-lg border border-border bg-background px-2 text-xs ${choice.enabled ? '' : 'opacity-60'}`}>
      <Checkbox checked={checked} disabled={!choice.enabled} onCheckedChange={() => toggleProvider(choice.providerKey, selectedKeys, onChange)} />
      <span>{choice.displayName}</span>
      {!choice.configured && <Badge variant="outline">未配置</Badge>}
      {choice.configured && !choice.enabled && <Badge variant="secondary">停用</Badge>}
    </label>
  );
}

function toggleProvider(providerKey: string, selectedKeys: string[], onChange: (keys: string[]) => void) {
  onChange(selectedKeys.includes(providerKey) ? selectedKeys.filter((key) => key !== providerKey) : [...selectedKeys, providerKey]);
}

function canSearch(applicationKey: string, countryISO2: string, callingCode: string, providerKeys: string[]) {
  return applicationKey.trim() !== '' && (countryISO2.trim() !== '' || callingCode.trim() !== '') && providerKeys.length > 0;
}

function SearchHint({ applicationKey, countryISO2, countryCallingCode, providerKeys }: { applicationKey: string; countryISO2: string; countryCallingCode: string; providerKeys: string[] }) {
  const hints = [];
  if (!applicationKey.trim()) hints.push('填写应用');
  if (!countryISO2.trim() && !countryCallingCode.trim()) hints.push('填写国家 ISO2 或区号');
  if (providerKeys.length === 0) hints.push('至少选择一个启用平台');
  if (hints.length === 0) return <p className="text-xs text-muted-foreground">支持刷新和分享当前查询链接；后端只查询当前选中的接码平台。</p>;
  return <p className="text-xs text-muted-foreground">还需要：{hints.join('、')}</p>;
}

function numberInputValue(value: string | number | null) {
  const number = Number(value || 0);
  return Number.isFinite(number) ? number : 0;
}

function compareQueryFromSearch(params: URLSearchParams): CompareQuery {
  return { applicationKey: params.get('application_key') || '', countryISO2: params.get('country_iso2') || '', countryCallingCode: params.get('country_calling_code') || '', providerKeys: params.getAll('provider_key').filter(Boolean), minAvailable: Math.max(0, numberInputValue(params.get('min_available') || '1')), sort: offerSort(params.get('sort')) };
}

function smsPriceOfferQuery(query: CompareQuery): SmsPriceOfferQuery {
  return { applicationKey: query.applicationKey, countryISO2: query.countryISO2, countryCallingCode: query.countryCallingCode, providerKeys: query.providerKeys, minAvailable: query.minAvailable };
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
