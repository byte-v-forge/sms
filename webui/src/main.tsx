import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App as AntApp, ConfigProvider } from 'antd';
import 'antd/dist/reset.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider } from 'react-router/dom';
import { smsRouter } from './dashboard/sms-routes';
import './styles.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 10_000,
      refetchOnWindowFocus: false
    }
  }
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ConfigProvider theme={{ token: { borderRadius: 10, colorPrimary: '#2563eb' } }}>
      <AntApp>
        <QueryClientProvider client={queryClient}>
          <RouterProvider router={smsRouter} />
        </QueryClientProvider>
      </AntApp>
    </ConfigProvider>
  </StrictMode>
);
