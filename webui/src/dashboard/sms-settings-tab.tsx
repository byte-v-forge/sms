import { useEffect, useState } from 'react';
import { LoaderCircle, Save, Trash2 } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { Badge, Button, Card, CardContent, CardDescription, CardHeader, CardTitle, DescriptionLine, EmptyBlock, Input, Switch } from '../ui';
import { getSmsProviderBalance, smsKeys, type SaveSmsProviderSettingRequest, type SmsProviderOption, type SmsProviderSetting } from './sms-api';
import { moneyText } from './sms-format';

type ProviderCardDescriptor = {
  option: SmsProviderOption;
  config?: SmsProviderSetting;
};

type SmsSettingsTabProps = {
  providerOptions: SmsProviderOption[];
  configs: SmsProviderSetting[];
  busy?: boolean;
  savingProviderKey?: string;
  deletingProviderKey?: string;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
};

export function SmsSettingsTab(props: SmsSettingsTabProps) {
  const providers = mergeProviderOptions(props.providerOptions, props.configs);
  const enabledCount = props.configs.filter((config) => config.enabled).length;
  return (
    <div className="min-h-0 flex-1 overflow-auto bg-muted/20 p-4">
      <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
        <div><div className="text-sm font-semibold">SMS Provider</div><div className="text-xs text-muted-foreground">每个平台在卡片内直接配置，配置仅保存在 SMS 服务。</div></div>
        <div className="flex flex-wrap items-center gap-2"><Badge variant="secondary">{props.providerOptions.length} 插件</Badge><Badge variant="outline">{props.configs.length} 已配置</Badge><Badge variant={enabledCount > 0 ? 'default' : 'secondary'}>{enabledCount} 启用</Badge></div>
      </div>
      <div className="flex flex-wrap items-start gap-3">
        {providers.map(({ option, config }) => <ProviderCard key={option.provider_key} config={config} provider={option} saving={props.savingProviderKey === option.provider_key} deleting={props.deletingProviderKey === option.provider_key} onSave={props.onSave} onDelete={props.onDelete} />)}
        {props.busy && providers.length === 0 && <LoaderCircle className="size-5 animate-spin text-muted-foreground" />}
        {!props.busy && providers.length === 0 && <EmptyBlock text="暂无接码源插件" />}
      </div>
    </div>
  );
}

function ProviderCard({ config, provider, saving, deleting, onSave, onDelete }: {
  config?: SmsProviderSetting;
  provider: SmsProviderOption;
  saving?: boolean;
  deleting?: boolean;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
}) {
  const [apiKey, setAPIKey] = useState('');
  const [enabled, setEnabled] = useState(true);
  useEffect(() => {
    setAPIKey('');
    setEnabled(config?.enabled ?? true);
  }, [config?.enabled, config?.provider_key]);
  const balance = useQuery({ queryKey: smsKeys.balance(provider.provider_key), queryFn: () => getSmsProviderBalance(provider.provider_key), enabled: !!config?.provider_key && config.enabled && !!config.api_key_set, refetchInterval: 60000 });
  const canSave = enabled === false || !!apiKey.trim() || !!config?.api_key_set;
  return (
    <Card className="w-[340px] max-w-full flex-none bg-background">
      <CardHeader>
        <CardTitle className="flex items-center justify-between gap-2"><span className="truncate">{provider.display_name || provider.provider_key}</span>{providerStatusBadge(config, enabled)}</CardTitle>
        <CardDescription>{provider.provider_key}</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-3">
        <Input type="password" placeholder={config?.api_key_set ? '留空保留 API Key' : 'Provider API Key'} value={apiKey} onChange={(event) => setAPIKey(event.target.value)} />
        <div className="flex items-center justify-between rounded-lg border border-border p-2 text-sm"><span>启用接码源</span><Switch checked={enabled} onCheckedChange={setEnabled} /></div>
        <div className="grid gap-1 rounded-lg bg-muted/30 p-2"><DescriptionLine label="余额" value={balanceText(config, balance.isLoading, moneyText(balance.data?.balance))} /><DescriptionLine label="API Key" value={config?.api_key_set ? '已配置' : '未配置'} /></div>
        <div className="flex justify-end gap-2">
          <Button aria-label="保存" title="保存" size="icon" disabled={saving || !canSave} onClick={() => onSave({ provider_key: provider.provider_key, enabled, api_key: apiKey.trim() || undefined })}>{saving ? <LoaderCircle className="size-4 animate-spin" /> : <Save className="size-4" />}</Button>
          <Button aria-label="删除" title="删除" variant="outline" size="icon" disabled={!config?.provider_key || deleting} onClick={() => onDelete(provider.provider_key)}>{deleting ? <LoaderCircle className="size-4 animate-spin" /> : <Trash2 className="size-4" />}</Button>
        </div>
      </CardContent>
    </Card>
  );
}

function mergeProviderOptions(providerOptions: SmsProviderOption[], configs: SmsProviderSetting[]): ProviderCardDescriptor[] {
  const configsByKey = new Map(configs.map((config) => [config.provider_key, config]));
  const providers = providerOptions.map((option) => ({ option, config: configsByKey.get(option.provider_key) }));
  const optionKeys = new Set(providerOptions.map((option) => option.provider_key));
  for (const config of configs) if (!optionKeys.has(config.provider_key)) providers.push({ option: { provider_key: config.provider_key, display_name: config.provider_key }, config });
  return providers;
}

function balanceText(config: SmsProviderSetting | undefined, loading: boolean, balance: string) {
  if (!config) return '-';
  return loading ? '读取中' : balance;
}

function providerStatusBadge(config: SmsProviderSetting | undefined, enabled: boolean) {
  if (!config) return <Badge variant="outline">未配置</Badge>;
  return enabled ? <Badge>启用</Badge> : <Badge variant="secondary">停用</Badge>;
}
