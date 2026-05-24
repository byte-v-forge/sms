import { useEffect, useMemo, useState } from 'react';
import { MessageSquareText } from 'lucide-react';
import { Tabs, TabsContent, TabsList, TabsTrigger, ToastMessage, WorkspaceToolbar, useMutation, useQuery, useQueryClient, useToastMessage } from '@/dashboard/module-kit';
import type { SmsProviderConfig, SmsRouteProfile, UpsertProviderConfigResponse, UpsertRouteProfileResponse } from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';
import { cancelSmsActivation, deleteSmsProviderConfig, deleteSmsRouteProfile, listSmsActivations, listSmsProviderConfigs, listSmsProviderPlugins, listSmsRouteProfiles, saveSmsProviderConfig, saveSmsRouteProfile, smsKeys } from './sms-api';
import { newSmsProviderConfig, newSmsRouteProfile } from './sms-format';
import { OrdersTab } from './orders-tab';
import { ProviderTab } from './provider-tab';
import { RouteProfilesTab } from './route-profiles-tab';

export function SmsPage() {
  const queryClient = useQueryClient();
  const toast = useToastMessage();
  const [selectedConfigId, setSelectedConfigId] = useState('');
  const [selectedProfileKey, setSelectedProfileKey] = useState('');
  const pluginsQuery = useQuery({ queryKey: smsKeys.plugins, queryFn: listSmsProviderPlugins });
  const configsQuery = useQuery({ queryKey: smsKeys.configs, queryFn: listSmsProviderConfigs });
  const profilesQuery = useQuery({ queryKey: smsKeys.profiles, queryFn: listSmsRouteProfiles });
  const activationsQuery = useQuery({ queryKey: smsKeys.activations, queryFn: listSmsActivations, refetchInterval: 5000 });
  const configs = configsQuery.data?.configs || [];
  const plugins = pluginsQuery.data?.plugins || [];
  const profiles = profilesQuery.data?.profiles || [];
  const selectedConfig = useMemo(() => configs.find((item: SmsProviderConfig) => item.provider_config_id === selectedConfigId) || null, [configs, selectedConfigId]);
  const selectedProfile = useMemo(() => profiles.find((item: SmsRouteProfile) => item.profile_key === selectedProfileKey) || null, [profiles, selectedProfileKey]);

  useEffect(() => {
    if (!selectedConfigId && configs[0]?.provider_config_id) setSelectedConfigId(configs[0].provider_config_id);
  }, [configs, selectedConfigId]);
  useEffect(() => {
    if (!selectedProfileKey && profiles[0]?.profile_key) setSelectedProfileKey(profiles[0].profile_key);
  }, [profiles, selectedProfileKey]);

  const saveMutation = useMutation({
    mutationFn: saveSmsProviderConfig,
    onSuccess: async (resp: UpsertProviderConfigResponse) => {
      if (resp.config?.provider_config_id) setSelectedConfigId(resp.config.provider_config_id);
      await queryClient.invalidateQueries({ queryKey: smsKeys.configs });
      toast.showOK('Provider配置已保存');
    },
    onError: toast.showError
  });
  const deleteMutation = useMutation({
    mutationFn: deleteSmsProviderConfig,
    onSuccess: async () => {
      setSelectedConfigId('');
      await queryClient.invalidateQueries({ queryKey: smsKeys.configs });
      toast.showOK('Provider配置已删除');
    },
    onError: toast.showError
  });
  const saveProfileMutation = useMutation({
    mutationFn: saveSmsRouteProfile,
    onSuccess: async (resp: UpsertRouteProfileResponse) => {
      if (resp.profile?.profile_key) setSelectedProfileKey(resp.profile.profile_key);
      await queryClient.invalidateQueries({ queryKey: smsKeys.profiles });
      toast.showOK('Profile配置已保存');
    },
    onError: toast.showError
  });
  const deleteProfileMutation = useMutation({
    mutationFn: deleteSmsRouteProfile,
    onSuccess: async () => {
      setSelectedProfileKey('');
      await queryClient.invalidateQueries({ queryKey: smsKeys.profiles });
      toast.showOK('Profile配置已删除');
    },
    onError: toast.showError
  });
  const cancelMutation = useMutation({
    mutationFn: cancelSmsActivation,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: smsKeys.activations });
      toast.showOK('号码已取消');
    },
    onError: toast.showError
  });

  return (
    <>
      <ToastMessage toast={toast.toast} />
      <section className="workspace singlePaneWorkspace">
        <div className="panel">
          <Tabs defaultValue="profiles" className="flex min-h-0 flex-1 flex-col">
            <WorkspaceToolbar
              title={<span className="inline-flex items-center gap-2"><MessageSquareText className="size-4" />SMS</span>}
              meta={`${profiles.length}个Profile · ${configs.length}个Provider · ${activationsQuery.data?.activations?.length || 0}个订单`}
              tabs={<TabsList><TabsTrigger value="profiles">Profile</TabsTrigger><TabsTrigger value="providers">Provider</TabsTrigger><TabsTrigger value="orders">号码订单</TabsTrigger></TabsList>}
            />
            <TabsContent value="profiles" className="mt-0 min-h-0 flex-1">
              <RouteProfilesTab
                plugins={plugins}
                profiles={profiles}
                selected={selectedProfile || (selectedProfileKey === 'new' ? newSmsRouteProfile() : null)}
                busy={profilesQuery.isLoading}
                saving={saveProfileMutation.isPending}
                deleting={deleteProfileMutation.isPending}
                onSelect={setSelectedProfileKey}
                onNew={() => setSelectedProfileKey('new')}
                onSave={(profile: SmsRouteProfile) => saveProfileMutation.mutate(profile)}
                onDelete={(key) => deleteProfileMutation.mutate(key)}
              />
            </TabsContent>
            <TabsContent value="providers" className="mt-0 min-h-0 flex-1">
              <ProviderTab
                plugins={plugins}
                configs={configs}
                selected={selectedConfig || (selectedConfigId === 'new' ? newSmsProviderConfig() : null)}
                busy={configsQuery.isLoading}
                saving={saveMutation.isPending}
                deleting={deleteMutation.isPending}
                onSelect={setSelectedConfigId}
                onNew={() => setSelectedConfigId('new')}
                onSave={(config: SmsProviderConfig) => saveMutation.mutate(config)}
                onDelete={(id) => deleteMutation.mutate(id)}
              />
            </TabsContent>
            <TabsContent value="orders" className="mt-0 min-h-0 flex-1">
              <OrdersTab
                activations={activationsQuery.data?.activations || []}
                cancelingId={cancelMutation.variables}
                onCancel={(id) => cancelMutation.mutate(id)}
              />
            </TabsContent>
          </Tabs>
        </div>
      </section>
    </>
  );
}
