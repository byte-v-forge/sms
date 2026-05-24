import { Plus } from 'lucide-react';
import { Badge, Button, Card, CardContent, CardHeader, CardTitle, useQuery } from '@/dashboard/module-kit';
import type { SmsProviderConfig, SmsProviderPluginDescriptor } from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';
import { getSmsProviderBalance, smsKeys } from './sms-api';
import { moneyText, newSmsProviderConfig } from './sms-format';
import { ProviderConfigForm } from './provider-config-form';

type ProviderTabProps = {
  plugins: SmsProviderPluginDescriptor[];
  configs: SmsProviderConfig[];
  selected: SmsProviderConfig | null;
  busy?: boolean;
  saving?: boolean;
  deleting?: boolean;
  onSelect: (id: string) => void;
  onNew: () => void;
  onSave: (config: SmsProviderConfig) => void;
  onDelete: (id: string) => void;
};

export function ProviderTab(props: ProviderTabProps) {
  return (
    <div className="grid min-h-0 flex-1 grid-cols-[minmax(0,1fr)_360px]">
      <div className="min-h-0 overflow-auto p-3">
        <div className="mb-3 flex items-center justify-between">
          <div className="text-sm font-semibold">Provider</div>
          <Button size="sm" onClick={props.onNew}><Plus className="size-4" />新增</Button>
        </div>
        <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-3">
          {props.configs.map((config) => (
            <ProviderCard
              key={config.provider_config_id}
              config={config}
              plugin={props.plugins.find((item) => item.provider_key === config.provider_key)}
              selected={props.selected?.provider_config_id === config.provider_config_id}
              onSelect={() => props.onSelect(config.provider_config_id)}
            />
          ))}
          {!props.busy && props.configs.length === 0 && <div className="empty">暂无Provider配置</div>}
        </div>
      </div>
      <ProviderConfigForm
        config={props.selected || newSmsProviderConfig()}
        plugins={props.plugins}
        saving={props.saving}
        deleting={props.deleting}
        onSave={props.onSave}
        onDelete={props.onDelete}
      />
    </div>
  );
}

function ProviderCard({ config, plugin, selected, onSelect }: { config: SmsProviderConfig; plugin?: SmsProviderPluginDescriptor; selected: boolean; onSelect: () => void }) {
  const balance = useQuery({
    queryKey: smsKeys.balance(config.provider_config_id),
    queryFn: () => getSmsProviderBalance(config.provider_config_id),
    enabled: !!config.provider_config_id && config.enabled && config.credential_secret_set,
    refetchInterval: 60000
  });
  return (
    <Card className={selected ? 'border-primary' : ''} role="button" tabIndex={0} onClick={onSelect}>
      <CardHeader className="space-y-1 p-3">
        <div className="flex items-center justify-between gap-2">
          <CardTitle className="truncate text-sm">{config.display_name || config.provider_config_id}</CardTitle>
          <Badge variant={config.enabled ? 'default' : 'secondary'}>{config.enabled ? '启用' : '停用'}</Badge>
        </div>
        <div className="truncate text-xs text-muted-foreground">{plugin?.display_name || config.provider_key} · {config.provider_config_id}</div>
      </CardHeader>
      <CardContent className="grid gap-1 p-3 pt-0 text-xs">
        <Line label="余额" value={balance.isLoading ? '读取中' : moneyText(balance.data?.balance)} />
        <Line label="密钥" value={config.credential_secret_set ? '已配置' : '未配置'} />
      </CardContent>
    </Card>
  );
}

function Line({ label, value }: { label: string; value: string }) {
  return <div className="flex min-w-0 justify-between gap-3"><span className="text-muted-foreground">{label}</span><span className="truncate font-medium">{value}</span></div>;
}
