import { useEffect, useMemo, useState } from 'react';
import { MessageSquareText } from 'lucide-react';
import { createHotStreamURL, ToastMessage, useHotStreamInvalidation, WorkspaceTabbedPanel, useMutation, useQuery, useQueryClient, useToastMessage } from '@byte-v-forge/common-ui';
import { cancelSmsOrder, deleteSmsProviderSetting, listSmsOrders, listSmsProviderSettings, saveSmsProviderSetting, smsKeys, type SaveSmsProviderSettingResponse, type SmsProviderSetting } from './sms-api';
import { OrdersTab } from './orders-tab';
import { SmsSettingsTab } from './sms-settings-tab';

export function SmsPage() {
  const queryClient = useQueryClient();
  const toast = useToastMessage();
  const [selectedProviderKey, setSelectedProviderKey] = useState('');
  const settingsQuery = useQuery({ queryKey: smsKeys.settingsProviders, queryFn: listSmsProviderSettings });
  const ordersQuery = useQuery({ queryKey: smsKeys.orders, queryFn: listSmsOrders });
  const configs = settingsQuery.data?.providers || [];
  const options = settingsQuery.data?.provider_options || [];
  const selectedConfig = useMemo(() => configs.find((item: SmsProviderSetting) => item.provider_key === selectedProviderKey) || null, [configs, selectedProviderKey]);

  useEffect(() => {
    if (!selectedProviderKey && configs[0]?.provider_key) setSelectedProviderKey(configs[0].provider_key);
  }, [configs, selectedProviderKey]);
  useHotStreamInvalidation({
    url: createHotStreamURL('/api/sms', { eventTypes: ['sms.order.updated', 'sms.provider_config.updated', 'sms.provider_config.deleted'] }),
    rules: [
      { queryKey: smsKeys.orders, eventTypes: ['sms.order.updated'], resourceTypes: ['sms.order'] },
      { queryKey: smsKeys.settingsProviders, eventTypes: ['sms.provider_config.updated', 'sms.provider_config.deleted'], resourceTypes: ['sms.provider_config'] }
    ]
  });

  const saveMutation = useMutation({
    mutationFn: saveSmsProviderSetting,
    onSuccess: async (resp: SaveSmsProviderSettingResponse) => {
      if (resp.provider?.provider_key) setSelectedProviderKey(resp.provider.provider_key);
      await queryClient.invalidateQueries({ queryKey: smsKeys.settingsProviders });
      toast.showOK('接码源已保存');
    },
    onError: toast.showError
  });
  const deleteMutation = useMutation({
    mutationFn: deleteSmsProviderSetting,
    onSuccess: async () => {
      setSelectedProviderKey('');
      await queryClient.invalidateQueries({ queryKey: smsKeys.settingsProviders });
      toast.showOK('接码源已删除');
    },
    onError: toast.showError
  });
  const cancelMutation = useMutation({
    mutationFn: cancelSmsOrder,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: smsKeys.orders });
      toast.showOK('号码已取消');
    },
    onError: toast.showError
  });

  return (
    <>
      <ToastMessage toast={toast.toast} />
      <WorkspaceTabbedPanel
        defaultValue="orders"
        title={<span className="inline-flex items-center gap-2"><MessageSquareText className="size-4" />SMS</span>}
        meta={`${configs.length}个接码源 · ${ordersQuery.data?.orders?.length || 0}个订单`}
        tabs={[
          {
            value: 'orders',
            label: '号码订单',
            content: <OrdersTab orders={ordersQuery.data?.orders || []} cancelingId={cancelMutation.variables} onCancel={(id) => cancelMutation.mutate(id)} />
          },
          {
            value: 'settings',
            label: '设置',
            content: <SmsSettingsTab providerOptions={options} configs={configs} selected={selectedConfig} busy={settingsQuery.isLoading} saving={saveMutation.isPending} deleting={deleteMutation.isPending} onSelect={setSelectedProviderKey} onNew={() => setSelectedProviderKey('new')} onSave={(input) => saveMutation.mutate(input)} onDelete={(id) => deleteMutation.mutate(id)} />
          }
        ]}
      />
    </>
  );
}
