const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8081';
const API_KEY = import.meta.env.VITE_ADMIN_API_KEY || '';

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init?.headers as Record<string, string>),
  };

  // Add API key for admin endpoints
  if (API_KEY && path.startsWith('/v1/admin')) {
    headers['X-API-Key'] = API_KEY;
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