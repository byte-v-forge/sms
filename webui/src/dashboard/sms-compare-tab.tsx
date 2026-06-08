import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { SearchOutlined } from '@ant-design/icons';
import { Button, Card, Checkbox, Flex, Form, Input, InputNumber, Select, Space, Tag, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { listSmsPriceOffers, smsKeys, type SmsPriceOfferQuery, type SmsProviderOption, type SmsProviderSetting } from './sms-api';
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
  providerOptions: SmsProviderOption[];
  configs: SmsProviderSetting[];
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
    setSearchParams(compareQuerySearchParams(draftQuery(sort)));
  }

  function changeSort(next: OfferSort) {
    setSort(next);
    const draft = draftQuery(next);
    if (canSearch(draft.applicationKey, draft.countryISO2, draft.countryCallingCode, draft.providerKeys)) setSearchParams(compareQuerySearchParams(draft));
  }

  function draftQuery(nextSort: OfferSort): CompareQuery {
    return {
      applicationKey: applicationKey.trim(),
      countryISO2: countryISO2.trim().toUpperCase(),
      countryCallingCode: countryCallingCode.trim().replace(/^\+/, ''),
      providerKeys: [...activeKeys],
      minAvailable: Math.max(0, minAvailable),
      sort: nextSort
    };
  }

  return (
    <Flex className="sms-fill" vertical>
      <Card size="small" style={{ margin: 16, marginBottom: 0 }}>
        <Form layout="vertical" onSubmitCapture={submitQuery}>
          <Flex gap={12} wrap="wrap" align="end">
            <Form.Item label="应用" style={{ minWidth: 210, flex: 1, marginBottom: 8 }}><Input placeholder="whatsapp/gojek" value={applicationKey} onChange={(event) => setApplicationKey(event.target.value)} /></Form.Item>
            <Form.Item label="国家 ISO2" style={{ minWidth: 140, marginBottom: 8 }}><Input placeholder="ID" value={countryISO2} onChange={(event) => setCountryISO2(event.target.value)} /></Form.Item>
            <Form.Item label="国家区号" style={{ minWidth: 140, marginBottom: 8 }}><Input placeholder="62" value={countryCallingCode} onChange={(event) => setCountryCallingCode(event.target.value)} /></Form.Item>
            <Form.Item label="最低库存" style={{ minWidth: 140, marginBottom: 8 }}><InputNumber min={0} value={minAvailable} onChange={(value) => setMinAvailable(numberInputValue(value))} /></Form.Item>
            <Form.Item label="排序" style={{ minWidth: 140, marginBottom: 8 }}><Select value={sort} onChange={changeSort} options={[{ label: '按低价', value: 'price' }, { label: '按库存', value: 'available' }, { label: '按平台', value: 'provider' }]} /></Form.Item>
            <Button type="primary" htmlType="submit" icon={<SearchOutlined />} disabled={!canSearch(applicationKey, countryISO2, countryCallingCode, activeKeys)} style={{ marginBottom: 8 }}>查询比对</Button>
          </Flex>
          <ProviderPicker choices={choices} selectedKeys={activeKeys} onChange={setSelectedKeys} />
        </Form>
      </Card>
      <CompareSummary loading={offersQuery.isLoading} total={offers.length} providerCount={new Set(offers.map((offer) => offer.provider_key)).size} best={top} error={error} />
      <OffersTable offers={offers} top={top} loading={offersQuery.isLoading} queried={queried} error={error} acquiringOfferId={acquiringOfferId} onAcquire={onAcquire} />
    </Flex>
  );
}

function ProviderPicker({ choices, selectedKeys, onChange }: { choices: ProviderChoice[]; selectedKeys: string[]; onChange: (keys: string[]) => void }) {
  if (choices.length === 0) return <Typography.Text type="secondary">暂无 provider 插件</Typography.Text>;
  return (
    <Space size={[12, 8]} wrap>
      {choices.map((choice) => (
        <Checkbox key={choice.providerKey} checked={selectedKeys.includes(choice.providerKey)} disabled={!choice.enabled} onChange={() => toggleProvider(choice.providerKey, selectedKeys, onChange)}>
          {choice.displayName} {!choice.configured && <Tag>未配置</Tag>}
        </Checkbox>
      ))}
    </Space>
  );
}

function toggleProvider(providerKey: string, selectedKeys: string[], onChange: (keys: string[]) => void) {
  onChange(selectedKeys.includes(providerKey) ? selectedKeys.filter((key) => key !== providerKey) : [...selectedKeys, providerKey]);
}

function canSearch(applicationKey: string, countryISO2: string, callingCode: string, providerKeys: string[]) {
  return applicationKey.trim() !== '' && (countryISO2.trim() !== '' || callingCode.trim() !== '') && providerKeys.length > 0;
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
