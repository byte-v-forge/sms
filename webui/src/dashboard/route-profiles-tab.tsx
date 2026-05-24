import { Plus } from 'lucide-react';
import { Badge, Button, Card, CardContent, CardHeader, CardTitle } from '@/dashboard/module-kit';
import type { SmsProviderPluginDescriptor, SmsRouteProfile } from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';
import { newSmsRouteProfile, strategyText } from './sms-format';
import { RouteProfileForm } from './route-profile-form';

type Props = {
  plugins: SmsProviderPluginDescriptor[];
  profiles: SmsRouteProfile[];
  selected: SmsRouteProfile | null;
  busy?: boolean;
  saving?: boolean;
  deleting?: boolean;
  onSelect: (key: string) => void;
  onNew: () => void;
  onSave: (profile: SmsRouteProfile) => void;
  onDelete: (key: string) => void;
};

export function RouteProfilesTab(props: Props) {
  return (
    <div className="grid min-h-0 flex-1 grid-cols-[minmax(0,1fr)_520px]">
      <div className="min-h-0 overflow-auto p-3">
        <div className="mb-3 flex items-center justify-between">
          <div className="text-sm font-semibold">Profile</div>
          <Button size="sm" onClick={props.onNew}><Plus className="size-4" />新增</Button>
        </div>
        <div className="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-3">
          {props.profiles.map((profile) => (
            <ProfileCard
              key={profile.profile_key}
              profile={profile}
              selected={props.selected?.profile_key === profile.profile_key}
              onSelect={() => props.onSelect(profile.profile_key)}
            />
          ))}
          {!props.busy && props.profiles.length === 0 && <div className="empty">暂无Profile配置</div>}
        </div>
      </div>
      <RouteProfileForm
        plugins={props.plugins}
        profile={props.selected || newSmsRouteProfile()}
        saving={props.saving}
        deleting={props.deleting}
        onSave={props.onSave}
        onDelete={props.onDelete}
      />
    </div>
  );
}

function ProfileCard({ profile, selected, onSelect }: { profile: SmsRouteProfile; selected: boolean; onSelect: () => void }) {
  return (
    <Card className={selected ? 'border-primary' : ''} role="button" tabIndex={0} onClick={onSelect}>
      <CardHeader className="space-y-1 p-3">
        <div className="flex items-center justify-between gap-2">
          <CardTitle className="truncate text-sm">{profile.display_name || profile.profile_key}</CardTitle>
          <Badge variant={profile.enabled ? 'default' : 'secondary'}>{profile.enabled ? '启用' : '停用'}</Badge>
        </div>
        <div className="truncate text-xs text-muted-foreground">{profile.profile_key}</div>
      </CardHeader>
      <CardContent className="grid gap-1 p-3 pt-0 text-xs">
        <Line label="策略" value={strategyText(profile.selection_strategy)} />
        <Line label="Routes" value={String(profile.routes?.length || 0)} />
      </CardContent>
    </Card>
  );
}

function Line({ label, value }: { label: string; value: string }) {
  return <div className="flex min-w-0 justify-between gap-3"><span className="text-muted-foreground">{label}</span><span className="truncate font-medium">{value}</span></div>;
}
