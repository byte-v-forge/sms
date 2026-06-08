import { CheckSquare, Square } from 'lucide-react';
import { Badge, Button, Checkbox } from '../ui';
import type { ProviderChoice } from './sms-compare-data';

type ProviderPickerProps = {
  choices: ProviderChoice[];
  selectedKeys: string[];
  onChange: (keys: string[]) => void;
};

export function ProviderPicker({ choices, selectedKeys, onChange }: ProviderPickerProps) {
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
