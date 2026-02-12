const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8081';

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init?.headers as Record<string, string>),
  };

  // Add API key for admin endpoints from sessionStorage or env
  if (path.startsWith('/v1/admin')) {
    const key = sessionStorage.getItem('ml_admin_key')
      || import.meta.env.VITE_ADMIN_API_KEY
      || '';
    if (key) {
      headers['X-API-Key'] = key;
    }
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers,
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}