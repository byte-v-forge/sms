import { useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { ChevronDown, ChevronRight, Save, Trash2 } from 'lucide-react';
import { Button, Input, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/dashboard/module-kit';
import {
  SmsConfigFieldKind,
  SmsConfigFieldTarget,
  type SmsProviderConfig,
  type SmsProviderConfigField,
  type SmsProviderPluginDescriptor
} from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';
import { defaultSmsProviderPolicy, durationSeconds, newSmsProviderConfig, secondsDuration } from './sms-format';

type FormProps = {
  config: SmsProviderConfig | null;
  plugins: SmsProviderPluginDescriptor[];
  saving?: boolean;
  deleting?: boolean;
  onSave: (config: SmsProviderConfig) => void;
  onDelete: (id: string) => void;
};

export function ProviderConfigForm({ config, plugins, saving, deleting, onSave, onDelete }: FormProps) {
  const [draft, setDraft] = useState<SmsProviderConfig>(() => config || newSmsProviderConfig());
  const [advancedOpen, setAdvancedOpen] = useState(false);
  useEffect(() => setDraft(config || newSmsProviderConfig()), [config]);
  const providerType = draft.provider_key || '5sim';
  const plugin = plugins.find((item) => item.provider_key === providerType) || plugins[0];
  const fields = plugin?.config_fields || [];

  function patch(next: Partial<SmsProviderConfig>) {
    setDraft((current) => ({ ...current, ...next }));
  }

  function patchProviderType(value: string) {
    const nextPlugin = plugins.find((item) => item.provider_key === value);
    patch({ provider_key: value, provider_config_id: value, display_name: nextPlugin?.display_name || value, enabled: true, policy: defaultSmsProviderPolicy(value) });
  }

  function patchField(field: SmsProviderConfigField, value: string) {
    setDraft((current) => writeConfigField(current, field, value));
  }

  function patchPolicy(field: keyof NonNullable<SmsProviderConfig['policy']>, seconds: number) {
    patch({ policy: { ...(draft.policy || defaultSmsProviderPolicy(providerType)), [field]: secondsDuration(seconds) } });
  }

  function save() {
    onSave({
      ...draft,
      provider_config_id: providerType,
      provider_key: providerType,
      display_name: plugin?.display_name || providerType,
      enabled: true,
      api_endpoint: '',
      credential_secret_ref: '',
      proxy_ref: '',
      default_target: undefined,
      upstream_service_key: '',
      provider_country_id: '',
      policy: draft.policy || defaultSmsProviderPolicy(providerType),
      labels: normalizeLabels(draft.labels || {})
    });
  }
  const policy = draft.policy || defaultSmsProviderPolicy(providerType);
  const basicFields = fields.filter((field) => !field.advanced);
  const advancedFields = fields.filter((field) => field.advanced);

  return (
    <div className="flex min-h-0 flex-col gap-3 border-l border-border/70 p-3">
      <div className="text-sm font-semibold">Provider配置</div>
      <div className="grid gap-2">
        <Field label="Provider Type">
          <Select value={providerType} onValueChange={patchProviderType}>
            <SelectTrigger><SelectValue /></SelectTrigger>
            <SelectContent>
              {plugins.map((item) => <SelectItem key={item.provider_key} value={item.provider_key}>{item.display_name || item.provider_key}</SelectItem>)}
            </SelectContent>
          </Select>
        </Field>
        <ProviderFieldList fields={basicFields} draft={draft} onChange={patchField} />
        {advancedFields.length > 0 && (
          <div className="grid gap-2">
            <Button type="button" variant="ghost" size="sm" className="w-fit gap-1 px-1" onClick={() => setAdvancedOpen((open) => !open)}>
              {advancedOpen ? <ChevronDown className="size-4" /> : <ChevronRight className="size-4" />}高级配置
            </Button>
            {advancedOpen && <ProviderFieldList fields={advancedFields} draft={draft} onChange={patchField} />}
          </div>
        )}
        <div className="grid grid-cols-2 gap-2">
          <Field label="有效期(分钟)"><Input type="number" min={1} value={Math.round(durationSeconds(policy.activation_ttl) / 60)} onChange={(e) => patchPolicy('activation_ttl', Number(e.target.value) * 60)} /></Field>
          <Field label="轮询(秒)"><Input type="number" min={1} value={durationSeconds(policy.poll_interval)} onChange={(e) => patchPolicy('poll_interval', Number(e.target.value))} /></Field>
          <Field label="取消等待(秒)"><Input type="number" min={0} value={durationSeconds(policy.cancel_allowed_after)} onChange={(e) => patchPolicy('cancel_allowed_after', Number(e.target.value))} /></Field>
          <Field label="提前取消重试(秒)"><Input type="number" min={0} value={durationSeconds(policy.early_cancel_retry_after)} onChange={(e) => patchPolicy('early_cancel_retry_after', Number(e.target.value))} /></Field>
        </div>
      </div>
      <div className="mt-auto flex gap-2">
        <Button className="flex-1" disabled={saving} onClick={save}><Save className="size-4" />保存</Button>
        <Button variant="outline" size="icon" disabled={!draft.provider_config_id || deleting} onClick={() => onDelete(draft.provider_config_id)}><Trash2 className="size-4" /></Button>
      </div>
    </div>
  );
}

function ProviderFieldList({ fields, draft, onChange }: {
  fields: SmsProviderConfigField[];
  draft: SmsProviderConfig;
  onChange: (field: SmsProviderConfigField, value: string) => void;
}) {
  return <>{fields.map((field) => <ProviderField key={field.field_key} field={field} draft={draft} onChange={(value) => onChange(field, value)} />)}</>;
}

function ProviderField({ field, draft, onChange }: {
  field: SmsProviderConfigField;
  draft: SmsProviderConfig;
  onChange: (value: string) => void;
}) {
  const type = field.kind === SmsConfigFieldKind.SMS_CONFIG_FIELD_KIND_SECRET
    ? 'password'
    : field.kind === SmsConfigFieldKind.SMS_CONFIG_FIELD_KIND_NUMBER || field.kind === SmsConfigFieldKind.SMS_CONFIG_FIELD_KIND_DURATION_SECONDS ? 'number' : 'text';
  const placeholder = field.kind === SmsConfigFieldKind.SMS_CONFIG_FIELD_KIND_SECRET && draft.credential_secret_set ? '留空则保留现有密钥' : field.placeholder;
  return (
    <Field label={`${field.label}${field.required ? ' *' : ''}`}>
      <Input type={type} placeholder={placeholder} value={readConfigField(draft, field)} onChange={(e) => onChange(e.target.value)} />
    </Field>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <div className="grid gap-1"><Label>{label}</Label>{children}</div>;
}

function readConfigField(config: SmsProviderConfig, field: SmsProviderConfigField) {
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_CREDENTIAL_SECRET) return config.credential_secret || '';
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_API_ENDPOINT) return config.api_endpoint || '';
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_HTTP_PROXY) return config.http_proxy || '';
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_LABEL) return config.labels?.[field.field_key] || '';
  return '';
}

function writeConfigField(config: SmsProviderConfig, field: SmsProviderConfigField, value: string): SmsProviderConfig {
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_CREDENTIAL_SECRET) return { ...config, credential_secret: value };
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_API_ENDPOINT) return { ...config, api_endpoint: value };
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_HTTP_PROXY) return { ...config, http_proxy: value };
  if (field.target === SmsConfigFieldTarget.SMS_CONFIG_FIELD_TARGET_LABEL) {
    return { ...config, labels: { ...(config.labels || {}), [field.field_key]: value } };
  }
  return config;
}

function normalizeLabels(labels: Record<string, string>) {
  return Object.fromEntries(Object.entries(labels).filter(([key, value]) => key.trim() && String(value).trim()));
}
