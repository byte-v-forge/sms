import type { FormEvent, ReactNode } from 'react';
import { RotateCcw, Search } from 'lucide-react';
import { Button, Card, Input, Select } from '../ui';
import type { OfferSort, ProviderChoice } from './sms-compare-data';
import { numberInputValue } from './sms-compare-query';
import { ProviderPicker } from './sms-provider-picker';

type SmsCompareFormProps = {
  choices: ProviderChoice[];
  applicationKey: string;
  countryISO2: string;
  countryCallingCode: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
  canSubmit: boolean;
  onApplicationKeyChange: (value: string) => void;
  onCountryISO2Change: (value: string) => void;
  onCountryCallingCodeChange: (value: string) => void;
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
        <div className="grid gap-2 md:grid-cols-6">
          <Field label="应用" className="md:col-span-2"><Input placeholder="whatsapp/gojek" value={props.applicationKey} onChange={(event) => props.onApplicationKeyChange(event.target.value)} /></Field>
          <Field label="国家 ISO2"><Input placeholder="ID" value={props.countryISO2} onChange={(event) => props.onCountryISO2Change(event.target.value)} /></Field>
          <Field label="国家区号"><Input placeholder="62" value={props.countryCallingCode} onChange={(event) => props.onCountryCallingCodeChange(event.target.value)} /></Field>
          <Field label="最低库存"><Input min={0} type="number" value={props.minAvailable} onChange={(event) => props.onMinAvailableChange(numberInputValue(event.target.value))} /></Field>
          <Field label="排序"><Select value={props.sort} onChange={(event) => props.onSortChange(event.target.value as OfferSort)}><option value="price">按低价</option><option value="available">按库存</option><option value="provider">按平台</option></Select></Field>
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <ProviderPicker choices={props.choices} selectedKeys={props.providerKeys} onChange={props.onProviderKeysChange} />
          <div className="flex items-center gap-2">
            <Button aria-label="重置查询条件" title="重置查询条件" size="icon-sm" variant="outline" onClick={props.onReset}><RotateCcw className="size-4" /></Button>
            <Button disabled={!props.canSubmit} type="submit"><Search className="size-4" />查询比对</Button>
          </div>
        </div>
        <SearchHint applicationKey={props.applicationKey} countryISO2={props.countryISO2} countryCallingCode={props.countryCallingCode} providerKeys={props.providerKeys} />
      </form>
    </Card>
  );
}

function Field({ label, className, children }: { label: string; className?: string; children: ReactNode }) {
  return <label className={`grid gap-1 text-xs font-medium text-muted-foreground ${className || ''}`}><span>{label}</span>{children}</label>;
}

function SearchHint({ applicationKey, countryISO2, countryCallingCode, providerKeys }: { applicationKey: string; countryISO2: string; countryCallingCode: string; providerKeys: string[] }) {
  const hints = [];
  if (!applicationKey.trim()) hints.push('填写应用');
  if (!countryISO2.trim() && !countryCallingCode.trim()) hints.push('填写国家 ISO2 或区号');
  if (providerKeys.length === 0) hints.push('至少选择一个启用平台');
  if (hints.length === 0) return <p className="text-xs text-muted-foreground">支持刷新和分享当前查询链接；后端只查询当前选中的接码平台。</p>;
  return <p className="text-xs text-muted-foreground">还需要：{hints.join('、')}</p>;
}
