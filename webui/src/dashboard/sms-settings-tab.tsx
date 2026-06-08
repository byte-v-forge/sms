import { useEffect, useState } from 'react';
import { Save, Trash2 } from 'lucide-react';
import { Badge, Button, DescriptionLine, EmptyBlock, Input, Item, ItemContent, ItemDescription, ItemTitle, Switch, useQuery } from '../ui';
import { getSmsProviderBalance, smsKeys, type SaveSmsProviderSettingRequest, type SmsProviderOption, type SmsProviderSetting } from './sms-api';
import { moneyText } from './sms-format';

type ProviderCardDescriptor = {
  option: SmsProviderOption;
  config?: SmsProviderSetting;
};

type ProviderBadgeVariant = 'default' | 'secondary' | 'outline';

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
        <div>
          <div className="text-sm font-semibold">SMS Provider</div>
          <div className="text-xs text-muted-foreground">每个平台在卡片内直接配置。</div>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="secondary">{props.providerOptions.length} 插件</Badge>
          <Badge variant="outline">{props.configs.length} 已配置</Badge>
          <Badge variant={enabledCount > 0 ? 'default' : 'secondary'}>{enabledCount} 启用</Badge>
        </div>
      </div>
      <div className="flex flex-wrap items-start gap-3">
        {providers.map(({ option, config }) => (
          <ProviderCard
            key={option.provider_key}
            config={config}
            provider={option}
            saving={props.savingProviderKey === option.provider_key}
            deleting={props.deletingProviderKey === option.provider_key}
            onSave={props.onSave}
            onDelete={props.onDelete}
          />
        ))}
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
  const balance = useQuery({
    queryKey: smsKeys.balance(provider.provider_key),
    queryFn: () => getSmsProviderBalance(provider.provider_key),
    enabled: !!config?.provider_key && config.enabled && !!config.api_key_set,
    refetchInterval: 60000
  });
  const hasCredential = !!apiKey.trim() || !!config?.api_key_set;
  const canSave = enabled === false || hasCredential;
  const statusText = providerStatusText(config, enabled);
  return (
    <Item className="w-[320px] max-w-full flex-none bg-background transition hover:border-primary/40 hover:bg-muted/20" variant="outline">
      <ItemContent className="min-w-0 gap-3">
        <div>
          <ItemTitle className="w-full justify-between">
            <span className="truncate">{provider.display_name || provider.provider_key}</span>
            <Badge aria-label={statusText} className="h-5 w-5 rounded-full px-0" variant={providerBadgeVariant(config, enabled)}>
              <span aria-hidden="true" className="size-1.5 rounded-full bg-current" />
            </Badge>
          </ItemTitle>
          <ItemDescription>{provider.provider_key}</ItemDescription>
        </div>
        <div className="grid gap-2">
          <Input
            type="password"
            placeholder={config?.api_key_set ? '留空保留 API Key' : 'Provider API Key'}
            value={apiKey}
            onChange={(event) => setAPIKey(event.target.value)}
          />
          <div className="flex h-9 items-center justify-end rounded-lg border border-border/70 px-3">
            <Switch aria-label="启用" checked={enabled} onCheckedChange={setEnabled} />
          </div>
        </div>
        <div className="grid gap-1 rounded-lg bg-muted/30 p-2">
          <DescriptionLine label="余额" value={balanceText(config, balance.isLoading, moneyText(balance.data?.balance))} />
          <DescriptionLine label="API Key" value={config?.api_key_set ? '已配置' : '未配置'} />
        </div>
        <div className="flex justify-end gap-2">
          <Button
            aria-label="保存"
            size="icon"
            disabled={saving || !canSave}
            onClick={() => onSave({ provider_key: provider.provider_key, enabled, api_key: apiKey.trim() || undefined })}
          >
            <Save className="size-4" />
          </Button>
          <Button aria-label="删除" variant="outline" size="icon" disabled={!config?.provider_key || deleting} onClick={() => onDelete(provider.provider_key)}>
            <Trash2 className="size-4" />
          </Button>
        </div>
      </ItemContent>
    </Item>
  );
}

function mergeProviderOptions(providerOptions: SmsProviderOption[], configs: SmsProviderSetting[]): ProviderCardDescriptor[] {
  const configsByKey = new Map(configs.map((config) => [config.provider_key, config]));
  const providers = providerOptions.map((option) => ({ option, config: configsByKey.get(option.provider_key) }));
  const optionKeys = new Set(providerOptions.map((option) => option.provider_key));
  for (const config of configs) {
    if (!optionKeys.has(config.provider_key)) providers.push({ option: { provider_key: config.provider_key, display_name: config.provider_key }, config });
  }
  return providers;
}

function balanceText(config: SmsProviderSetting | undefined, loading: boolean, balance: string) {
  if (!config) return '-';
  return loading ? '读取中' : balance;
}

function providerStatusText(config: SmsProviderSetting | undefined, enabled: boolean) {
  if (!config) return '未配置';
  return enabled ? '启用' : '停用';
}

function providerBadgeVariant(config: SmsProviderSetting | undefined, enabled: boolean): ProviderBadgeVariant {
  if (!config) return 'outline';
  return enabled ? 'default' : 'secondary';
}
