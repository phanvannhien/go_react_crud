const API_BASE_URL = 'http://localhost:8080';

interface APIResponse<T> {
  data: T;
  error: null | {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

interface PaginatedResponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
  error: null | {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

interface CursorResponse<T> {
  data: T[];
  next_cursor: string;
  error: null | {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

class APIError extends Error {
  code: string;
  details?: Record<string, string>;

  constructor(code: string, message: string, details?: Record<string, string>) {
    super(message);
    this.code = code;
    this.details = details;
  }
}

function getToken(): string | null {
  return sessionStorage.getItem('auth_token');
}

function setToken(token: string): void {
  sessionStorage.setItem('auth_token', token);
}

function clearToken(): void {
  sessionStorage.removeItem('auth_token');
}

async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (res.status === 401) {
    clearToken();
    window.location.href = '/login';
    throw new APIError('unauthorized', 'Session expired');
  }

  const json = await res.json();

  if (json.error) {
    throw new APIError(
      json.error.code || 'unknown',
      json.error.message || 'An error occurred',
      json.error.details
    );
  }

  return json as T;
}

function buildQueryString(params: Record<string, string | number | boolean | undefined | null>): string {
  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      searchParams.set(key, String(value));
    }
  });
  const qs = searchParams.toString();
  return qs ? `?${qs}` : '';
}

export { apiFetch, buildQueryString, getToken, setToken, clearToken, APIError };
export type { APIResponse, PaginatedResponse, CursorResponse };
