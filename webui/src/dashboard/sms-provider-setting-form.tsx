import { useEffect, useState } from 'react';
import { Save, Trash2 } from 'lucide-react';
import { Button, DashboardField, Input, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@byte-v-forge/common-ui';
import type { SaveSmsProviderSettingRequest, SmsProviderOption, SmsProviderSetting } from './sms-api';

type ProviderDraft = SmsProviderSetting & {
  api_key?: string;
};

type FormProps = {
  config: SmsProviderSetting | null;
  providerOptions: SmsProviderOption[];
  saving?: boolean;
  deleting?: boolean;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
};

export function SmsProviderSettingForm({ config, providerOptions, saving, deleting, onSave, onDelete }: FormProps) {
  const defaultProvider = providerOptions[0];
  const [draft, setDraft] = useState<ProviderDraft>(() => config || newProviderDraft(defaultProvider));
  useEffect(() => setDraft(config || newProviderDraft(defaultProvider)), [config, defaultProvider]);
  const providerType = draft.provider_key || defaultProvider?.provider_key || '';
  const existingProviderKey = config?.provider_key || '';
  const hasAPIKey = draft.api_key_set || !!draft.api_key?.trim();
  const canSave = !!providerType && (draft.enabled === false || hasAPIKey);

  function patch(next: Partial<ProviderDraft>) {
    setDraft((current) => ({ ...current, ...next }));
  }

  function patchProviderType(value: string) {
    const nextProvider = providerOptions.find((item) => item.provider_key === value);
    setDraft(newProviderDraft(nextProvider));
  }

  function save() {
    onSave({
      provider_key: providerType,
      enabled: draft.enabled !== false,
      api_key: draft.api_key?.trim() || undefined
    });
  }

  return (
    <div className="flex min-h-0 flex-col gap-3 border-l border-border/70 p-3">
      <div>
        <div className="text-sm font-semibold">接码源接入</div>
        <div className="text-xs text-muted-foreground">只配置平台接入必需项，高级参数使用服务默认值。</div>
      </div>
      <div className="grid gap-2">
        <DashboardField label="平台">
          <Select value={providerType} onValueChange={patchProviderType}>
            <SelectTrigger><SelectValue /></SelectTrigger>
            <SelectContent>
              {providerOptions.map((item) => <SelectItem key={item.provider_key} value={item.provider_key}>{item.display_name || item.provider_key}</SelectItem>)}
            </SelectContent>
          </Select>
        </DashboardField>
        <DashboardField label="API Key *">
          <Input
            type="password"
            placeholder={draft.api_key_set ? '留空则保留已配置 API Key' : '输入 Provider API Key'}
            value={draft.api_key || ''}
            onChange={(e) => patch({ api_key: e.target.value })}
          />
        </DashboardField>
        <DashboardField label="状态">
          <Button type="button" variant={draft.enabled !== false ? 'default' : 'outline'} onClick={() => patch({ enabled: draft.enabled === false })}>
            {draft.enabled !== false ? '已启用' : '已停用'}
          </Button>
        </DashboardField>
      </div>
      <div className="mt-auto flex gap-2">
        <Button className="flex-1" disabled={saving || !canSave} onClick={save}><Save className="size-4" />保存</Button>
        <Button variant="outline" size="icon" disabled={!existingProviderKey || deleting} onClick={() => onDelete(existingProviderKey)}><Trash2 className="size-4" /></Button>
      </div>
    </div>
  );
}

function newProviderDraft(provider?: SmsProviderOption): ProviderDraft {
  const providerKey = provider?.provider_key || '';
  return {
    provider_key: providerKey,
    enabled: true,
    api_key_set: false,
    api_key: ''
  };
}
