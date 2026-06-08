import { useMemo } from 'react';
import { MessageSquareText } from 'lucide-react';
import { toast } from 'sonner';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Outlet } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsOrderCodeView, SmsOrderView, SmsProviderConfig, SmsProviderPluginDescriptor } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { createHotStreamURL, useHotStreamInvalidation, WorkspaceRoutedPanel } from '../ui';
import {
  acquireSmsFromOffer,
  cancelSmsOrder,
  deleteSmsProviderConfig,
  listSmsProviderConfigs,
  listSmsProviderPlugins,
  listSmsOrderCodes,
  listSmsOrders,
  saveSmsProviderConfig,
  smsKeys
} from './sms-api';

export type SmsPageContext = {
  providerOptions: SmsProviderPluginDescriptor[];
  configs: SmsProviderConfig[];
  settingsBusy: boolean;
  orders: SmsOrderView[];
  codes: SmsOrderCodeView[];
  cancelingId?: string;
  acquiringOfferId?: string;
  savingProviderKey?: string;
  deletingProviderKey?: string;
  onAcquire: (offer: SmsPriceOffer) => void;
  onCancel: (id: string) => void;
  onSave: (config: SmsProviderConfig) => void;
  onDelete: (id: string) => void;
};

export function SmsPage() {
  const queryClient = useQueryClient();
  const providerPluginsQuery = useQuery({ queryKey: smsKeys.providerPlugins, queryFn: listSmsProviderPlugins });
  const providerConfigsQuery = useQuery({ queryKey: smsKeys.providerConfigs, queryFn: listSmsProviderConfigs });
  const ordersQuery = useQuery({ queryKey: smsKeys.orders, queryFn: listSmsOrders });
  const orderIds = (ordersQuery.data?.orders || []).map((item) => item.order?.order_id || '').filter(Boolean);
  const codesQuery = useQuery({ queryKey: smsKeys.orderCodes(orderIds), queryFn: () => listSmsOrderCodes(orderIds, 5), enabled: orderIds.length > 0 });
  const configs = providerConfigsQuery.data?.configs || [];
  const options = providerPluginsQuery.data?.plugins || [];
  const streamRules = useMemo(() => [
    { queryKey: smsKeys.orders, eventTypes: ['sms.order.updated'], resourceTypes: ['sms.order'] },
    { queryKey: smsKeys.orderCodesRoot, eventTypes: ['sms.order.updated'], resourceTypes: ['sms.order'] },
    { queryKey: smsKeys.providerConfigs, eventTypes: ['sms.provider_config.updated', 'sms.provider_config.deleted'], resourceTypes: ['sms.provider_config'] }
  ], []);

  useHotStreamInvalidation({
    url: createHotStreamURL('/api/sms', { eventTypes: ['sms.order.updated', 'sms.provider_config.updated', 'sms.provider_config.deleted'] }),
    rules: streamRules
  });

  const saveMutation = useMutation({ mutationFn: saveSmsProviderConfig, onSuccess: async () => { await queryClient.invalidateQueries({ queryKey: smsKeys.providerConfigs }); toast.success('接码源已保存'); }, onError: showError });
  const deleteMutation = useMutation({ mutationFn: deleteSmsProviderConfig, onSuccess: async () => { await queryClient.invalidateQueries({ queryKey: smsKeys.providerConfigs }); toast.success('接码源已删除'); }, onError: showError });
  const cancelMutation = useMutation({ mutationFn: cancelSmsOrder, onSuccess: async () => { await queryClient.invalidateQueries({ queryKey: smsKeys.orders }); toast.success('号码已取消'); }, onError: showError });
  const acquireMutation = useMutation({
    mutationFn: acquireSmsFromOffer,
    onSuccess: async (response) => {
      if (response.error) {
        toast.error(response.error.message || response.error.code);
        return;
      }
      await queryClient.invalidateQueries({ queryKey: smsKeys.orders });
      toast.success('号码已获取');
    },
    onError: showError
  });
  const pageContext: SmsPageContext = {
    providerOptions: options,
    configs,
    settingsBusy: providerPluginsQuery.isLoading || providerConfigsQuery.isLoading,
    orders: ordersQuery.data?.orders || [],
    codes: codesQuery.data?.codes || [],
    cancelingId: cancelMutation.variables,
    acquiringOfferId: acquireMutation.variables?.offer_ref?.offer_id,
    savingProviderKey: saveMutation.variables?.provider_key,
    deletingProviderKey: deleteMutation.variables,
    onAcquire: (offer) => acquireMutation.mutate(offer),
    onCancel: (id) => cancelMutation.mutate(id),
    onSave: (input) => saveMutation.mutate(input),
    onDelete: (id) => deleteMutation.mutate(id)
  };

  return (
    <WorkspaceRoutedPanel
      title={<span className="inline-flex items-center gap-2"><MessageSquareText className="size-4" />SMS</span>}
      meta={`${configs.length}个接码源 · ${ordersQuery.data?.orders?.length || 0}个订单`}
      tabs={[{ to: '/compare', label: '平台比价' }, { to: '/orders', label: '号码订单' }, { to: '/settings', label: '设置' }]}
    >
      <Outlet context={pageContext} />
    </WorkspaceRoutedPanel>
  );
}

function showError(error: unknown) {
  toast.error(errorText(error));
}

function errorText(error: unknown) {
  if (error instanceof Error) return error.message;
  if (typeof error === 'string') return error;
  return '操作失败';
}
