import { useCallback, useState } from 'react';

type ToastState = {
  kind: 'ok' | 'error';
  message: string;
};

export function useToastMessage() {
  const [toast, setToast] = useState<ToastState>();
  const showOK = useCallback((message: string) => setToast({ kind: 'ok', message }), []);
  const showError = useCallback((error: unknown) => setToast({ kind: 'error', message: errorText(error) }), []);
  return { toast, showOK, showError };
}

export function ToastMessage({ toast }: { toast?: ToastState }) {
  if (!toast) return null;
  const color = toast.kind === 'ok' ? 'border-emerald-200 bg-emerald-50 text-emerald-700' : 'border-red-200 bg-red-50 text-red-700';
  return <div className={`fixed right-4 top-4 z-50 rounded-lg border px-4 py-2 text-sm shadow-lg ${color}`}>{toast.message}</div>;
}

function errorText(error: unknown) {
  if (error instanceof Error) return error.message;
  if (typeof error === 'string') return error;
  return '操作失败';
}
