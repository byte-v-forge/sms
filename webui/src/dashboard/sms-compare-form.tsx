import type { FormEvent, ReactNode } from 'react';
import { RotateCcw, Search } from 'lucide-react';
import { Button, Card, Input, Select } from '../ui';
import type { OfferSort, ProviderChoice } from './sms-compare-data';
import { numberInputValue } from './sms-compare-query';
import { ProviderPicker } from './sms-provider-picker';

type SmsCompareFormProps = {
  choices: ProviderChoice[];
  searchText: string;
  providerKeys: string[];
  minAvailable: number;
  sort: OfferSort;
  canSubmit: boolean;
  onSearchTextChange: (value: string) => void;
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
          <Field label="搜索报价" className="md:col-span-4"><Input placeholder="搜索应用、国家、区号或平台，例如 whatsapp id 62" value={props.searchText} onChange={(event) => props.onSearchTextChange(event.target.value)} /></Field>
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
  return <p className="text-xs text-muted-foreground">默认展示已启用平台的报价；搜索框只做过滤，不需要先填写应用或国家。</p>;
}
