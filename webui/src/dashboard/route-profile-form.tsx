import { useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { Plus, Save, Trash2 } from 'lucide-react';
import { Button, Input, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, useQuery } from '@/dashboard/module-kit';
import { SmsRouteSelectionStrategy, type SmsProviderPluginDescriptor, type SmsRouteCandidate, type SmsRouteProfile } from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';
import { newSmsRouteCandidate, newSmsRouteProfile } from './sms-format';
import { ProviderRouteFields } from './route-provider-adapters';
import { listSmsRouteOptions, smsKeys } from './sms-api';

type Props = {
  plugins: SmsProviderPluginDescriptor[];
  profile: SmsRouteProfile | null;
  saving?: boolean;
  deleting?: boolean;
  onSave: (profile: SmsRouteProfile) => void;
  onDelete: (key: string) => void;
};

export function RouteProfileForm({ plugins, profile, saving, deleting, onSave, onDelete }: Props) {
  const [draft, setDraft] = useState<SmsRouteProfile>(() => profile || newSmsRouteProfile());
  useEffect(() => setDraft(profile || newSmsRouteProfile()), [profile]);

  function patch(next: Partial<SmsRouteProfile>) {
    setDraft((current) => ({ ...current, ...next }));
  }
  function patchRoute(index: number, route: SmsRouteCandidate) {
    patch({ routes: draft.routes.map((item, i) => (i === index ? route : item)) });
  }
  function save() {
    onSave({ ...draft, display_name: draft.display_name || draft.profile_key, enabled: true });
  }

  return (
    <div className="flex min-h-0 flex-col gap-3 border-l border-border/70 p-3">
      <div className="text-sm font-semibold">Profile配置</div>
      <div className="grid grid-cols-2 gap-2">
        <Field label="Profile Key"><Input value={draft.profile_key} onChange={(e) => patch({ profile_key: e.target.value })} /></Field>
        <Field label="名称"><Input value={draft.display_name} onChange={(e) => patch({ display_name: e.target.value })} /></Field>
        <Field label="选择策略"><StrategySelect value={draft.selection_strategy} onChange={(selection_strategy) => patch({ selection_strategy })} /></Field>
        <Field label="指定Provider"><ProviderSelect plugins={plugins} value={draft.preferred_provider_key || plugins[0]?.provider_key || 'smsbower'} onChange={(preferred_provider_key) => patch({ preferred_provider_key })} /></Field>
      </div>
      <div className="flex items-center justify-between">
        <div className="text-sm font-medium">Routes</div>
        <Button size="sm" onClick={() => patch({ routes: [...draft.routes, newSmsRouteCandidate(plugins[0]?.provider_key)] })}><Plus className="size-4" />新增</Button>
      </div>
      <div className="min-h-0 overflow-auto">
        <div className="grid gap-3">
          {draft.routes.map((route, index) => (
            <RouteEditor key={index} plugins={plugins} route={route} onChange={(next) => patchRoute(index, next)} onRemove={() => patch({ routes: draft.routes.filter((_, i) => i !== index) })} />
          ))}
        </div>
      </div>
      <div className="mt-auto flex gap-2">
        <Button className="flex-1" disabled={saving} onClick={save}><Save className="size-4" />保存</Button>
        <Button variant="outline" size="icon" disabled={!draft.profile_key || deleting} onClick={() => onDelete(draft.profile_key)}><Trash2 className="size-4" /></Button>
      </div>
    </div>
  );
}

function RouteEditor({ plugins, route, onChange, onRemove }: { plugins: SmsProviderPluginDescriptor[]; route: SmsRouteCandidate; onChange: (route: SmsRouteCandidate) => void; onRemove: () => void }) {
  const patch = (next: Partial<SmsRouteCandidate>) => onChange({ ...route, ...next });
  const plugin = plugins.find((item) => item.provider_key === route.provider_key) || plugins[0];
  const options = useQuery({
    queryKey: smsKeys.routeOptions(route.provider_key || ''),
    queryFn: () => listSmsRouteOptions(route.provider_key),
    enabled: !!route.provider_key
  });
  return (
    <div className="grid gap-2 rounded-md border border-border/70 p-2">
      <div className="grid grid-cols-[1fr_120px_40px] gap-2">
        <Input placeholder="Route ID" value={route.route_id} onChange={(e) => patch({ route_id: e.target.value })} />
        <ProviderSelect plugins={plugins} value={route.provider_key || plugin?.provider_key || 'smsbower'} onChange={(provider_key) => onChange(newSmsRouteCandidate(provider_key))} />
        <Button variant="outline" size="icon" onClick={onRemove}><Trash2 className="size-4" /></Button>
      </div>
      <div className="grid grid-cols-4 gap-2">
        <Input placeholder="优先级" value={String(route.priority || '')} onChange={(e) => patch({ priority: Number(e.target.value) || 0 })} />
        <Input placeholder="应用" value={route.target?.application_key || ''} onChange={(e) => patchTarget(route, onChange, { application_key: e.target.value })} />
        <Input placeholder="国家" value={route.target?.country_iso2 || ''} onChange={(e) => patchTarget(route, onChange, { country_iso2: e.target.value })} />
        <Input placeholder="区号" value={route.target?.country_calling_code || ''} onChange={(e) => patchTarget(route, onChange, { country_calling_code: e.target.value })} />
      </div>
      <ProviderRouteFields route={route} fields={plugin?.route_fields || []} options={options.data?.options} onChange={onChange} />
    </div>
  );
}

function patchTarget(route: SmsRouteCandidate, onChange: (route: SmsRouteCandidate) => void, target: Record<string, string>) {
  onChange({ ...route, target: { application_key: '', country_iso2: '', country_calling_code: '', min_price: undefined, max_price: undefined, ...(route.target || {}), ...target } });
}

function ProviderSelect({ plugins, value, onChange }: { plugins: SmsProviderPluginDescriptor[]; value: string; onChange: (value: string) => void }) {
  return <Select value={value} onValueChange={onChange}><SelectTrigger><SelectValue /></SelectTrigger><SelectContent>{plugins.map((plugin) => <SelectItem key={plugin.provider_key} value={plugin.provider_key}>{plugin.display_name || plugin.provider_key}</SelectItem>)}</SelectContent></Select>;
}

function StrategySelect({ value, onChange }: { value?: SmsRouteSelectionStrategy; onChange: (value: SmsRouteSelectionStrategy) => void }) {
  return <Select value={value || SmsRouteSelectionStrategy.SMS_ROUTE_SELECTION_STRATEGY_PRIORITY} onValueChange={(v) => onChange(v as SmsRouteSelectionStrategy)}><SelectTrigger><SelectValue /></SelectTrigger><SelectContent><SelectItem value={SmsRouteSelectionStrategy.SMS_ROUTE_SELECTION_STRATEGY_PRIORITY}>按优先级</SelectItem><SelectItem value={SmsRouteSelectionStrategy.SMS_ROUTE_SELECTION_STRATEGY_LOWEST_PRICE}>最低价</SelectItem><SelectItem value={SmsRouteSelectionStrategy.SMS_ROUTE_SELECTION_STRATEGY_SPECIFIED_PROVIDER}>指定Provider</SelectItem></SelectContent></Select>;
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <div className="grid gap-1"><Label>{label}</Label>{children}</div>;
}
