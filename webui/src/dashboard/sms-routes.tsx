import { Navigate, createBrowserRouter, useOutletContext } from 'react-router';
import { OrdersTab } from './orders-tab';
import { SmsCompareTab } from './sms-compare-tab';
import { SmsPage, type SmsPageContext } from './sms-page';
import { SmsSettingsTab } from './sms-settings-tab';

export const smsRouter = createBrowserRouter([
  {
    path: '/',
    Component: SmsPage,
    children: [
      { index: true, element: <Navigate to="/compare" replace /> },
      { path: 'compare', Component: SmsCompareRoute },
      { path: 'orders', Component: SmsOrdersRoute },
      { path: 'settings', Component: SmsSettingsRoute },
      { path: '*', element: <Navigate to="/compare" replace /> }
    ]
  }
]);

function SmsCompareRoute() {
  const page = useSmsPage();
  return (
    <SmsCompareTab
      providerOptions={page.providerOptions}
      configs={page.configs}
      acquiringOfferId={page.acquiringOfferId}
      onAcquire={page.onAcquire}
    />
  );
}

function SmsOrdersRoute() {
  const page = useSmsPage();
  return <OrdersTab orders={page.orders} codes={page.codes} cancelingId={page.cancelingId} onCancel={page.onCancel} />;
}

function SmsSettingsRoute() {
  const page = useSmsPage();
  return (
    <SmsSettingsTab
      providerOptions={page.providerOptions}
      configs={page.configs}
      busy={page.settingsBusy}
      savingProviderKey={page.savingProviderKey}
      deletingProviderKey={page.deletingProviderKey}
      onSave={page.onSave}
      onDelete={page.onDelete}
    />
  );
}

function useSmsPage() {
  return useOutletContext<SmsPageContext>();
}
