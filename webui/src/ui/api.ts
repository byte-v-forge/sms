export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers || {})
    }
  });
  const text = await response.text();
  const data = text ? JSON.parse(text) : {};
  if (!response.ok) {
    throw new Error(errorMessage(data) || `${response.status} ${response.statusText}`);
  }
  return data as T;
}

function errorMessage(data: unknown) {
  if (!data || typeof data !== 'object') return '';
  const value = (data as { error?: unknown }).error;
  return typeof value === 'string' ? value : '';
}
