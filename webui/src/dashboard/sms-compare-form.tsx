import type { FormEvent, ReactNode } from 'react';
import { RotateCcw, Search } from 'lucide-react';
import { Button, Card, Input, Select } from '../ui';
import { type OfferSort, type ProviderChoice } from './sms-compare-data';
import { applicationChoiceLabel, countryChoiceLabel, countryChoiceValue, type ApplicationChoice, type CountryChoice } from './sms-compare-options';
import { numberInputValue } from './sms-compare-query';
import { ProviderPicker } from './sms-provider-picker';

type SmsCompareFormProps = {
  choices: ProviderChoice[];
  applications: ApplicationChoice[];
  countries: CountryChoice[];
  searchText: string;
  countryValue: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
  canSubmit: boolean;
  onSearchTextChange: (value: string) => void;
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
        <div className="grid gap-2 md:grid-cols-6">
          <Field label="应用服务" className="md:col-span-3">
            <Input autoComplete="off" list="sms-application-options" placeholder="搜索/选择服务，例如 WhatsApp" value={props.searchText} onChange={(event) => props.onSearchTextChange(event.target.value)} />
            <datalist id="sms-application-options">
              {props.applications.map((item) => <option key={item.applicationKey} value={applicationChoiceLabel(item)} />)}
            </datalist>
          </Field>
          <Field label="国家">
            <Select value={props.countryValue} onChange={(event) => props.onCountryChange(event.target.value)}>
              <option value="">全部国家</option>
              {props.countries.map((item) => <option key={countryChoiceValue(item)} value={countryChoiceValue(item)}>{countryChoiceLabel(item)}</option>)}
            </Select>
          </Field>
          <Field label="最低库存"><Input min={0} type="number" value={props.minAvailable} onChange={(event) => props.onMinAvailableChange(numberInputValue(event.target.value))} /></Field>
          <Field label="排序"><Select value={props.sort} onChange={(event) => props.onSortChange(event.target.value as OfferSort)}><option value="price">按低价</option><option value="available">按库存</option><option value="provider">按平台</option></Select></Field>
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <ProviderPicker choices={props.choices} selectedKeys={props.providerKeys} onChange={props.onProviderKeysChange} />
          <div className="flex items-center gap-2">
            <Button aria-label="重置查询条件" title="重置查询条件" size="icon-sm" variant="outline" onClick={props.onReset}><RotateCcw className="size-4" /></Button>
            <Button disabled={!props.canSubmit} type="submit"><Search className="size-4" />搜索</Button>
          </div>
        </div>
        <SearchHint providerKeys={props.providerKeys} />
      </form>
    </Card>
  );
}

function Field({ label, className, children }: { label: string; className?: string; children: ReactNode }) {
  return <label className={`grid gap-1 text-xs font-medium text-muted-foreground ${className || ''}`}><span>{label}</span>{children}</label>;
}

function SearchHint({ providerKeys }: { providerKeys: string[] }) {
  if (providerKeys.length === 0) return <p className="text-xs text-muted-foreground">至少启用并选择一个接码平台后才会加载报价。</p>;
  return <p className="text-xs text-muted-foreground">先从下拉里搜服务，再选国家；默认展示已启用平台的可用报价。</p>;
}
