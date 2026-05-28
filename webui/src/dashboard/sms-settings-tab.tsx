import { Plus } from 'lucide-react';
import { Badge, Button, DescriptionLine, EmptyBlock, Item, ItemContent, ItemDescription, ItemTitle, useQuery } from '@byte-v-forge/common-ui';
import { getSmsProviderBalance, smsKeys, type SaveSmsProviderSettingRequest, type SmsProviderOption, type SmsProviderSetting } from './sms-api';
import { moneyText } from './sms-format';
import { SmsProviderSettingForm } from './sms-provider-setting-form';

type SmsSettingsTabProps = {
  providerOptions: SmsProviderOption[];
  configs: SmsProviderSetting[];
  selected: SmsProviderSetting | null;
  busy?: boolean;
  saving?: boolean;
  deleting?: boolean;
  onSelect: (id: string) => void;
  onNew: () => void;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
};

export function SmsSettingsTab(props: SmsSettingsTabProps) {
  return (
    <div className="grid min-h-0 flex-1 grid-cols-[minmax(0,1fr)_360px]">
      <div className="min-h-0 overflow-auto p-3">
        <div className="mb-3 flex items-center justify-between">
          <div className="text-sm font-semibold">接码源</div>
          <Button size="sm" onClick={props.onNew}><Plus className="size-4" />新增</Button>
        </div>
        <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-3">
          {props.configs.map((config) => (
            <ProviderCard
              key={config.provider_key}
              config={config}
              provider={props.providerOptions.find((item) => item.provider_key === config.provider_key)}
              selected={props.selected?.provider_key === config.provider_key}
              onSelect={() => props.onSelect(config.provider_key)}
            />
          ))}
          {!props.busy && props.configs.length === 0 && <EmptyBlock text="暂无接码源配置" />}
        </div>
      </div>
      <SmsProviderSettingForm
        config={props.selected}
        providerOptions={props.providerOptions}
        saving={props.saving}
        deleting={props.deleting}
        onSave={props.onSave}
        onDelete={props.onDelete}
      />
    </div>
  );
}

function ProviderCard({ config, provider, selected, onSelect }: { config: SmsProviderSetting; provider?: SmsProviderOption; selected: boolean; onSelect: () => void }) {
  const balance = useQuery({
    queryKey: smsKeys.balance(config.provider_key),
    queryFn: () => getSmsProviderBalance(config.provider_key),
    enabled: !!config.provider_key && config.enabled && !!config.api_key_set,
    refetchInterval: 60000
  });
  return (
    <Item className={selected ? 'border-primary' : ''} variant="outline" role="button" tabIndex={0} onClick={onSelect}>
      <ItemContent className="min-w-0">
        <ItemTitle className="w-full justify-between">
          <span className="truncate">{provider?.display_name || config.provider_key}</span>
          <Badge variant={config.enabled ? 'default' : 'secondary'}>{config.enabled ? '启用' : '停用'}</Badge>
        </ItemTitle>
        <ItemDescription>{config.provider_key}</ItemDescription>
        <div className="grid gap-1 pt-1">
          <DescriptionLine label="余额" value={balance.isLoading ? '读取中' : moneyText(balance.data?.balance)} />
          <DescriptionLine label="API Key" value={config.api_key_set ? '已配置' : '未配置'} />
        </div>
      </ItemContent>
    </Item>
  );
}
