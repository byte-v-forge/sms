import { MessageSquareText } from 'lucide-react';
import { Outlet } from 'react-router';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import type { SmsOrderCodeView, SmsOrderView } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import {
  createHotStreamURL,
  ToastMessage,
  useHotStreamInvalidation,
  WorkspaceRoutedPanel,
  useMutation,
  useQuery,
  useQueryClient,
  useToastMessage
} from '../ui';
import {
  acquireSmsFromOffer,
  cancelSmsOrder,
  deleteSmsProviderSetting,
  listSmsOrderCodes,
  listSmsOrders,
  listSmsProviderSettings,
  type SaveSmsProviderSettingRequest,
  type SmsProviderOption,
  type SmsProviderSetting,
  saveSmsProviderSetting,
  smsKeys
} from './sms-api';

export type SmsPageContext = {
  providerOptions: SmsProviderOption[];
  configs: SmsProviderSetting[];
  settingsBusy: boolean;
  orders: SmsOrderView[];
  codes: SmsOrderCodeView[];
  cancelingId?: string;
  acquiringOfferId?: string;
  savingProviderKey?: string;
  deletingProviderKey?: string;
  onAcquire: (offer: SmsPriceOffer) => void;
  onCancel: (id: string) => void;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
};

export function SmsPage() {
  const queryClient = useQueryClient();
  const toast = useToastMessage();
  const settingsQuery = useQuery({ queryKey: smsKeys.settingsProviders, queryFn: listSmsProviderSettings });
  const ordersQuery = useQuery({ queryKey: smsKeys.orders, queryFn: listSmsOrders });
  const orderIds = (ordersQuery.data?.orders || []).map((item) => item.order?.order_id || '').filter(Boolean);
  const codesQuery = useQuery({
    queryKey: smsKeys.orderCodes(orderIds),
    queryFn: () => listSmsOrderCodes(orderIds, 5),
    enabled: orderIds.length > 0
  });
  const configs = settingsQuery.data?.providers || [];
  const options = settingsQuery.data?.provider_options || [];

  useHotStreamInvalidation({
    url: createHotStreamURL('/api/sms', { eventTypes: ['sms.order.updated', 'sms.provider_config.updated', 'sms.provider_config.deleted'] }),
    rules: [
      { queryKey: smsKeys.orders, eventTypes: ['sms.order.updated'], resourceTypes: ['sms.order'] },
      { queryKey: smsKeys.orderCodesRoot, eventTypes: ['sms.order.updated'], resourceTypes: ['sms.order'] },
      { queryKey: smsKeys.settingsProviders, eventTypes: ['sms.provider_config.updated', 'sms.provider_config.deleted'], resourceTypes: ['sms.provider_config'] }
    ]
  });

  const saveMutation = useMutation({
    mutationFn: saveSmsProviderSetting,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: smsKeys.settingsProviders });
      toast.showOK('接码源已保存');
    },
    onError: toast.showError
  });
  const deleteMutation = useMutation({
    mutationFn: deleteSmsProviderSetting,
    onSuccess: async () => {
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
  const acquireMutation = useMutation({
    mutationFn: acquireSmsFromOffer,
    onSuccess: async (response) => {
      if (response.error) {
        toast.showError(response.error.message || response.error.code);
        return;
      }
      await queryClient.invalidateQueries({ queryKey: smsKeys.orders });
      toast.showOK('号码已获取');
    },
    onError: toast.showError
  });
  const pageContext: SmsPageContext = {
    providerOptions: options,
    configs,
    settingsBusy: settingsQuery.isLoading,
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
    <>
      <ToastMessage toast={toast.toast} />
      <WorkspaceRoutedPanel
        title={<span className="inline-flex items-center gap-2"><MessageSquareText className="size-4" />SMS</span>}
        meta={`${configs.length}个接码源 · ${ordersQuery.data?.orders?.length || 0}个订单`}
        tabs={[
          { to: '/compare', label: '平台比价' },
          { to: '/orders', label: '号码订单' },
          { to: '/settings', label: '设置' }
        ]}
      >
        <Outlet context={pageContext} />
      </WorkspaceRoutedPanel>
    </>
  );
}
