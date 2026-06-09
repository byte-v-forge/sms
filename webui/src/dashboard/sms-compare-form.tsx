import type { FormEvent, ReactNode } from 'react';
import { RotateCcw, Search } from 'lucide-react';
import { Button, Card, Input, Select } from '../ui';
import { type OfferSort, type ProviderChoice } from './sms-compare-data';
import { applicationSelectOptions, countrySelectOptions, type ApplicationChoice, type CountryChoice } from './sms-compare-options';
import { numberInputValue } from './sms-compare-query';
import { ProviderPicker } from './sms-provider-picker';
import { SearchSelect } from './sms-search-select';

type SmsCompareFormProps = {
  choices: ProviderChoice[];
  applications: ApplicationChoice[];
  countries: CountryChoice[];
  serviceText: string;
  applicationId: string;
  countryValue: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
  canSubmit: boolean;
  countriesLoading: boolean;
  serviceSelected: boolean;
  onServiceTextChange: (value: string) => void;
  onApplicationChange: (value: string) => void;
  onCountryChange: (value: string) => void;
  onProviderKeysChange: (keys: string[]) => void;
  onMinAvailableChange: (value: number) => void;
  onSortChange: (value: OfferSort) => void;
  onSubmit: (event?: FormEvent) => void;
  onReset: () => void;
};

export function SmsCompareForm(props: SmsCompareFormProps) {
  return (
    <Card className="m-4 mb-0 p-3">
      <form className="grid gap-3" onSubmit={props.onSubmit}>
        <div className="grid gap-2 lg:grid-cols-12">
          <div className="lg:col-span-5"><SearchSelect label="1. 应用/服务" placeholder="搜索并选择服务名称" emptyText="输入服务名称后选择匹配项" value={props.applicationId} searchValue={props.serviceText} options={applicationSelectOptions(props.applications)} shouldFilter={false} contentClassName="lg:w-[32rem]" onSearchChange={props.onServiceTextChange} onValueChange={props.onApplicationChange} /></div>
          <div className="lg:col-span-3"><SearchSelect disabled={!props.serviceSelected} label="2. 支持国家" placeholder={props.serviceSelected ? "选择支持国家" : "先选择应用/服务"} emptyText={countryEmptyText(props.serviceSelected, props.countriesLoading)} value={props.countryValue} options={countrySelectOptions(props.countries)} contentClassName="w-[min(18rem,calc(100vw-2rem))]" onValueChange={props.onCountryChange} /></div>
          <Field label="最低库存"><Input min={0} type="number" value={props.minAvailable} onChange={(event) => props.onMinAvailableChange(numberInputValue(event.target.value))} /></Field>
          <Field label="排序" className="lg:col-span-2"><Select value={props.sort} onChange={(event) => props.onSortChange(event.target.value as OfferSort)}><option value="price">按低价</option><option value="available">按库存</option><option value="provider">按平台</option></Select></Field>
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <ProviderPicker choices={props.choices} selectedKeys={props.providerKeys} onChange={props.onProviderKeysChange} />
          <div className="flex items-center gap-2">
            <Button aria-label="重置查询条件" title="重置查询条件" size="icon-sm" variant="outline" onClick={props.onReset}><RotateCcw className="size-4" /></Button>
            <Button disabled={!props.canSubmit} type="submit"><Search className="size-4" />搜索</Button>
          </div>
        </div>
        <SearchHint providerKeys={props.providerKeys} serviceSelected={props.serviceSelected} />
      </form>
    </Card>
  );
}

function Field({ label, className, children }: { label: string; className?: string; children: ReactNode }) {
  return <label className={`grid gap-1 text-xs font-medium text-muted-foreground ${className || ''}`}><span>{label}</span>{children}</label>;
}

function SearchHint({ providerKeys, serviceSelected }: { providerKeys: string[]; serviceSelected: boolean }) {
  if (providerKeys.length === 0) return <p className="text-xs text-muted-foreground">至少启用并选择一个接码平台后才会加载服务和国家。</p>;
  if (!serviceSelected) return <p className="text-xs text-muted-foreground">先搜索并选择应用/服务；选中后再加载该服务在各平台支持的国家。</p>;
  return <p className="text-xs text-muted-foreground">国家下拉只展示当前服务可用国家；平台内部服务代码只用于后台路由，不在页面展示。</p>;
}

function countryEmptyText(serviceSelected: boolean, loading: boolean) {
  if (!serviceSelected) return '先选择应用/服务';
  return loading ? '正在加载支持国家' : '当前服务没有可用国家';
}
